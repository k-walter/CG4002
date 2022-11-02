package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

type State = [cmn.NPlayer]cmn.PlayerState

type Engine struct {
	// public
	state     State
	chEvent   chan *pb.Event  // from relay
	chEval    chan *pb.State  // from eval
	chGrenade chan *EGrenaded // from viz

	// private
	chShoot     chan *eShootTimeout
	ShieldErrNs uint64
	ShootErrNs  uint64
	running     bool
}

type IEvent interface {
	updateEngine(*Engine) bool // return self if updated
	alertVizEvent() *pb.Event
	updateVizState() bool
	updateEvalState() bool
}

func Make(a *cmn.Arg) *Engine {
	e := Engine{
		state:       State{cmn.NewState(), cmn.NewState()},
		chEvent:     make(chan *pb.Event, cmn.ChSz),
		chEval:      make(chan *pb.State, cmn.ChSz),
		chShoot:     make(chan *eShootTimeout, cmn.ChSz),
		chGrenade:   make(chan *EGrenaded, cmn.ChSz),
		ShieldErrNs: a.ShieldErrNs,
		ShootErrNs:  a.ShootErrNs,
		running:     true,
	}

	// Subscribe to channels
	cmn.SubOld(cmn.Event2Eng, func(i interface{}) {
		go func(i *pb.Event) { e.chEvent <- i }(i.(*pb.Event))
	})
	cmn.SubOld(cmn.State2Eng, func(i interface{}) {
		go func(i *pb.State) { e.chEval <- i }(i.(*pb.State))
	})
	cmn.SubOld(cmn.Grenade2Eng, func(i interface{}) {
		go func(i *EGrenaded) { e.chGrenade <- i }(i.(*EGrenaded))
	})

	return &e
}

func (e *Engine) Run() {
	for e.running {
		ev := e.waitAnyEvent() // Serial order
		if ev == nil {
			continue
		}

		if didUpdate := ev.updateEngine(e); !didUpdate {
			continue
		}

		if event2Viz := ev.alertVizEvent(); event2Viz != nil {
			cmn.PubOld(cmn.Event2Viz, event2Viz)
		}
		if ev.updateVizState() {
			cmn.PubOld(cmn.State2Viz, snapshot(e.state))
		}
		if ev.updateEvalState() {
			cmn.PubOld(cmn.State2Eval, snapshot(e.state))
		}

		// Set back to none if sent state
		if ev.updateVizState() || ev.updateEvalState() {
			resetAction(e.state)
		}
	}
}

func resetAction(state State) {
	state[0].Action = pb.Action_none
	state[1].Action = pb.Action_none
}

func (e *Engine) waitAnyEvent() IEvent {
	// WARNING if WA, order by timestamp in channels
	// OPTIMIZE combine viz event and state update
	select {
	case timeout := <-e.chShoot: // shoot miss
		return timeout
	case grenaded := <-e.chGrenade: // viz checks if grenade hit
		return grenaded
	case state := <-e.chEval: // eval updates state
		return &eEvalResp{State: state}
	case event := <-e.chEvent: // relay/infer actions
		switch event.Action {
		case pb.Action_none:
			return nil
		case pb.Action_shoot:
			return &eShoot{Event: event}
		case pb.Action_shot:
			return &eShot{Event: event}
		case pb.Action_grenade:
			return &eGrenade{Event: event}
		case pb.Action_reload:
			return &eReload{Event: event}
		case pb.Action_shield:
			return &eShield{Event: event}
		case pb.Action_shieldAvailable:
			return &eUnshield{Event: event}
		case pb.Action_logout:
			return &eLogout{Event: event}
		default:
			log.Fatal("Unknown action", event.Action)
		}
	}
	return nil
}

func (e *Engine) Close() {
	// TODO Close channels
}

// Get <from, to>
func (e *Engine) getStates(player uint32) (*cmn.PlayerState, *cmn.PlayerState) {
	switch player {
	case 1:
		return &e.state[0], &e.state[1]
	case 2:
		return &e.state[1], &e.state[0]
	default:
		log.Fatal("Unknown player ", player)
	}
	return nil, nil
}

// Returns true if died (and revived)
func inflict(u uint32, player *cmn.PlayerState, dmg uint32) {
	if player.Hp+player.ShieldHealth <= dmg {
		// Die & revive
		player.NumDeaths += 1
		player.Hp = cmn.HpMax
		player.ShieldHealth = 0
		player.Grenades = cmn.GrenadeMax
		player.NumShield = cmn.ShieldMax
		player.Bullets = cmn.BulletMax

		// RULE shield cooldown rest
		if player.ShieldExpireNs != cmn.ShieldRst {
			// Why async? Current event can be sent < sending unshield event to viz
			cmn.PubOld(cmn.Event2Eng, &pb.Event{
				Player: u,
				Time:   player.ShieldExpireNs,
				Action: pb.Action_shieldAvailable,
			})
			player.ShieldExpireNs = cmn.ShieldRst

		}

	} else if player.ShieldHealth <= dmg {
		// Dmg shield+health
		player.Hp -= dmg - player.ShieldHealth
		player.ShieldHealth = 0
	} else {
		// Dmg shield
		player.ShieldHealth -= dmg
	}
}

// Deep copy, to prevent race condition when sending + modifying in engine
func snapshot(state State) *pb.State {
	now := uint64(time.Now().UnixNano())
	snapPlayer := func(src *cmn.PlayerState, dst *pb.PlayerState) {
		// Deep copy bytes
		bin, err := proto.Marshal(src.PlayerState)
		if err != nil {
			log.Fatal(err)
		}
		err = proto.Unmarshal(bin, dst)
		if err != nil {
			log.Fatal(err)
		}

		// Set last shield time
		if src.ShieldExpireNs <= now { // expired?
			dst.ShieldTime = 0
		} else { // ticking
			// remaining time
			dst.ShieldTime = float64(src.ShieldExpireNs-now) / 1e9 // in seconds
		}
	}

	ans := pb.State{
		P1: &pb.PlayerState{},
		P2: &pb.PlayerState{},
	}
	snapPlayer(&state[0], ans.P1)
	snapPlayer(&state[1], ans.P2)
	return &ans
}

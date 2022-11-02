package engine

import (
	cmn "cg4002/eComm/common"
	"cg4002/eComm/eval"
	pb "cg4002/protos"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

type State = [cmn.NPlayer]cmn.PlayerState

type Engine struct {
	state   [2]PlayerImpl
	running bool
	rnd     cmn.RoundT
	eval    eval.IEval
	rtt     time.Duration

	// Channels
	chEvent   chan *pb.Event     // from relay, pynq
	chGrenade chan *pb.InFovResp // from viz
}

type PlayerImpl struct {
	*pb.PlayerState
	fsm PlayerFSM

	shieldExpiry  time.Time
	shieldTimeout *time.Timer

	shoot        map[uint32]time.Time
	shot         map[uint32]time.Time
	shootTimeout *time.Timer
}

func NewPlayer() PlayerImpl {
	return PlayerImpl{
		PlayerState:   cmn.NewPlayerState(),
		fsm:           waiting,
		shieldExpiry:  cmn.GameTime,
		shieldTimeout: time.NewTimer(0),
		shoot:         make(map[uint32]time.Time),
		shot:          make(map[uint32]time.Time),
		shootTimeout:  time.NewTimer(0),
	}
}

type PlayerFSM uint8

const (
	waiting PlayerFSM = iota
	threwGrenade
	shotBullet
	done
)

func Make(a *cmn.Arg) *Engine {
	e := Engine{
		state:   [2]PlayerImpl{NewPlayer(), NewPlayer()},
		running: true,
		rnd:     1,
		eval:    eval.MakeMock(a), // eval.Make(a)
		rtt:     10 * time.Millisecond,

		chEvent:   cmn.Sub[*pb.Event](cmn.EEvent),
		chGrenade: cmn.Sub[*pb.InFovResp](cmn.EInFov),
	}

	return &e
}

func (e *Engine) Run() {
	// Drain timers
	for _, s := range e.state {
		cmn.Drain(s.shieldTimeout.C)
		cmn.Drain(s.shootTimeout.C)
	}

	for e.running {
		select {
		case ev := <-e.chEvent:
			e.handleEvent(ev)
			e.sendEval() // set shield timeout
		case <-e.state[0].shootTimeout.C:
			handleShootTimeout(&e.state[0])
			e.sendEval()
		case <-e.state[1].shootTimeout.C:
			handleShootTimeout(&e.state[1])
			e.sendEval()

		case rsp := <-e.chGrenade: // request must be sent from above
			e.handleFov(rsp)
		case <-e.state[0].shieldTimeout.C:
			handleShieldAvail(&e.state[0])
		case <-e.state[1].shieldTimeout.C:
			handleShieldAvail(&e.state[1])
		}
	}
}

func handleShootTimeout(p *PlayerImpl) {
	if p.fsm != shotBullet {
		return
	}
	p.fsm = done
	if !p.shootTimeout.Stop() {
		cmn.Drain(p.shootTimeout.C)
	}
}

func (e *Engine) Close() {
	// TODO Close channels
}

// Returns if event modified state
func (e *Engine) handleEvent(ev *pb.Event) {
	// previous round?
	if cmn.RoundT(ev.Rnd) < e.rnd {
		return
	}

	u, v := e.GetPlayers(ev.Player)
	doneWithAction := func() {
		u.Action = ev.Action
		u.fsm = done
	}

	// change tmp state
	switch ev.Action {
	case pb.Action_grenade:
		if u.fsm != waiting {
			return
		}
		cmn.Pub(cmn.EEvent, &pb.Event{
			Player: ev.Player ^ 0b11,
			Time:   ev.Time,
			Rnd:    ev.Rnd,
			Action: pb.Action_checkFov,
		})
		u.fsm = threwGrenade
		return

	case pb.Action_checkFov: // from eng to viz
		return

	case pb.Action_reload:
		if u.fsm != waiting {
			return
		}
		doneWithAction()
		if u.Bullets > 0 {
			u.Bullets = cmn.BulletMax
		}
		return

	case pb.Action_logout:
		if u.fsm != waiting {
			return
		}
		doneWithAction()
		return

	case pb.Action_shield:
		if u.fsm != waiting {
			return
		}
		doneWithAction()
		// Shield running?
		if u.shieldExpiry.After(cmn.GameTime) {
			return
		}
		// Set state only, set timeout after eval resp
		u.ShieldHealth = cmn.ShieldHpMax
		u.ShieldTime = cmn.ShieldTime.Seconds()

	case pb.Action_shoot:
		if u.fsm != waiting {
			return
		}
		doneWithAction()

		// No bullets?
		if u.Bullets == 0 {
			return
		}
		u.Bullets -= 1

		// Match with shot
		if _, fnd := v.shoot[ev.ShootID]; fnd {
			log.Fatalf("should not have duplicate shootID %v\n", ev.ShootID)
		}
		v.shoot[ev.ShootID] = cmn.NsToTime(ev.Time)

		if matchShot(ev.ShootID, v) {
			inflict(v, cmn.BulletDmg)
		} else {
			// Set timeout
			u.fsm = shotBullet
			u.shootTimeout.Reset(cmn.ShootErr)
		}

	case pb.Action_shot:
		// Timeout or no bullets?
		if v.fsm == done {
			return
		}

		// Match with shoot
		if _, fnd := u.shot[ev.ShootID]; fnd {
			log.Fatalf("should not have duplicate shootID %v\n", ev.ShootID)
		}
		u.shot[ev.ShootID] = cmn.NsToTime(ev.Time)
		if !matchShot(ev.ShootID, u) {
			return
		}

		inflict(u, cmn.BulletDmg)
		cmn.EXIT_UNLESS(v.fsm == shotBullet)
		v.fsm = done
		if !v.shootTimeout.Stop() {
			cmn.Drain(v.shootTimeout.C)
		}

	case pb.Action_none:
		return

	case pb.Action_grenaded:
		fallthrough
	case pb.Action_shieldAvailable:
		log.Fatal("Invalid player state", ev.Action)
	default:
		log.Fatal("Unhandled action", ev.Action)
	}

	return
}

// Returns true if died (and revived)
func inflict(p *PlayerImpl, dmg uint32) {
	if p.Hp+p.ShieldHealth <= dmg {
		// Die & revive
		p.NumDeaths += 1
		p.Hp = cmn.HpMax
		p.ShieldHealth = 0
		p.Grenades = cmn.GrenadeMax
		p.NumShield = cmn.ShieldMax
		p.Bullets = cmn.BulletMax

		// RULE reset shield cooldown
		if p.shieldExpiry.After(cmn.GameTime) {
			handleShieldAvail(p)
		}

	} else if p.ShieldHealth <= dmg {
		// Dmg shield+health
		p.Hp -= dmg - p.ShieldHealth
		p.ShieldHealth = 0
	} else {
		// Dmg shield
		p.ShieldHealth -= dmg
	}
}

func handleShieldAvail(p *PlayerImpl) {
	p.ShieldHealth = 0
	p.shieldExpiry = cmn.GameTime
	if !p.shieldTimeout.Stop() {
		cmn.Drain(p.shieldTimeout.C)
	}
}

func (e *Engine) handleFov(rsp *pb.InFovResp) {
	// Previous round?
	if cmn.RoundT(rsp.Rnd) < e.rnd {
		return
	}

	// Did opponent throw?
	me, opp := e.GetPlayers(rsp.Player)
	if opp.fsm != threwGrenade {
		return
	}
	opp.fsm = done

	// Damage me
	if rsp.InFov {
		inflict(me, cmn.GrenadeDmg)
	}
}

func (e *Engine) GetPlayers(i uint32) (*PlayerImpl, *PlayerImpl) {
	switch i {
	case 1:
		return &e.state[0], &e.state[1]
	case 2:
		return &e.state[1], &e.state[0]
	default:
		log.Fatal("Unknown player ", i)
	}
	return nil, nil
}

func (e *Engine) sendEval() {
	// All done?
	for _, s := range e.state {
		if s.fsm != done {
			return
		}
	}

	// Snapshot state
	st := time.Now()
	s := pb.State{
		P1: snapshotPlayer(&e.state[0], st, e.rtt),
		P2: snapshotPlayer(&e.state[1], st, e.rtt),
	}

	// Tx + rx from eval
	e.eval.BlockingSend(&s)
	t := e.eval.BlockingRecv()

	// Update rtt
	en := time.Now()
	rtt := (1-cmn.LPF)*e.rtt.Seconds() + cmn.LPF*en.Sub(st).Seconds()
	e.rtt = time.Duration(rtt) * time.Second

	// Reset players
	e.diffPlayer(&e.state[0], t.P1, en)
	e.diffPlayer(&e.state[1], t.P2, en)

	// Publish events
	e.rnd += 1
	cmn.Pub(cmn.EEvalResp, &cmn.EvalResp{State: t, Time: st})
	cmn.Pub(cmn.ERound, e.rnd)
}

func (e *Engine) diffPlayer(u *PlayerImpl, v *pb.PlayerState, now time.Time) {
	switch v.Action {
	case pb.Action_shield:
		// Did throw?
		if v.ShieldTime > 0 {
			u.shieldExpiry = now.Add(time.Duration(v.ShieldTime) * time.Second).Add(-e.rtt / 2)
			if !u.shieldTimeout.Reset(u.shieldExpiry.Sub(time.Now())) {
				cmn.Drain(u.shieldTimeout.C)
			}
		}

	case pb.Action_grenade:
		return
	case pb.Action_none:
		return
	case pb.Action_reload:
		return
	case pb.Action_shoot:
		return // OPTIMIZE clear opp's shoot/shot
	case pb.Action_logout:
		e.running = false

	case pb.Action_shot:
		fallthrough
	case pb.Action_grenaded:
		fallthrough
	case pb.Action_shieldAvailable:
		fallthrough
	case pb.Action_checkFov:
		log.Fatal("Invalid state", v.Action)

	default:
		log.Fatal("Unhandled state", v.Action)
	}

	// Copy state
	u.PlayerState = proto.Clone(v).(*pb.PlayerState)

	// Reset state
	u.fsm = waiting
	if !u.shootTimeout.Stop() {
		cmn.Drain(u.shootTimeout.C)
	}
}

func snapshotPlayer(p *PlayerImpl, t time.Time, rtt time.Duration) *pb.PlayerState {
	// Running shield?
	if p.shieldExpiry.After(cmn.GameTime) {
		p.ShieldTime = (p.shieldExpiry.Sub(t) - (rtt / 2)).Seconds()
	}
	return p.PlayerState
}

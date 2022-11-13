package engine

import (
	cmn "cg4002/eComm/common"
	"cg4002/eComm/eval"
	pb "cg4002/protos"
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
	a       *cmn.Arg

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

func NewPlayer(a *cmn.Arg) PlayerImpl {
	return PlayerImpl{
		PlayerState:   cmn.NewPlayerState(a),
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
		state:   [2]PlayerImpl{NewPlayer(a), NewPlayer(a)},
		running: true,
		rnd:     1,
		eval:    nil,
		rtt:     10 * time.Millisecond,
		a:       a,

		chEvent:   cmn.Sub[*pb.Event](cmn.EEvent),
		chGrenade: cmn.Sub[*pb.InFovResp](cmn.EInFov),
	}

	// Use mock eval server
	if a.MockEval {
		e.eval = eval.MakeMock(a)
	} else {
		e.eval = eval.Make(a)
	}

	return &e
}

func (e *Engine) Run() {
	// Drain timers
	for _, s := range e.state {
		cmn.Drain(s.shieldTimeout.C)
		cmn.Drain(s.shootTimeout.C)
	}

	// Send initial round
	cmn.Pub(cmn.ERound, e.rnd)

	// Send shoot/shot on match for latency debugging
	defer log.Println("P1 shot", e.state[0].shot, "shoot", e.state[0].shoot)
	defer log.Println("P2 shot", e.state[1].shot, "shoot", e.state[1].shoot)

	for e.a.MockEval || e.running {
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
			e.sendEval()

		case <-e.state[0].shieldTimeout.C:
			handleShieldAvail(&e.state[0])
		case <-e.state[1].shieldTimeout.C:
			handleShieldAvail(&e.state[1])
		}
	}
}

func (e *Engine) Close() {
	// TODO Close channels
}

func (e *Engine) handleEvent(ev *pb.Event) {
	// previous round?
	if cmn.RoundT(ev.Rnd) < e.rnd {
		return
	}

	u, v := e.GetPlayers(ev.Player)
	log.Printf("Eng|Handling %v, state=%v\n", ev, u.fsm)
	doneWithAction := func() {
		u.Action = ev.Action
		u.fsm = done
		log.Printf("Eng|Player %v done\n", ev.Player)
	}

	// change tmp state
	switch ev.Action {
	case pb.Action_grenade:
		if u.fsm != waiting {
			return
		}
		doneWithAction()
		if u.Grenades == 0 {
			return
		}
		u.Grenades--

		cmn.Pub(cmn.EEvent, &pb.Event{
			Player: ev.Player,
			Time:   ev.Time,
			Rnd:    ev.Rnd,
			Action: pb.Action_checkFov,
		})
		u.fsm = threwGrenade

	case pb.Action_checkFov: // from eng to viz

	case pb.Action_reload:
		if u.fsm != waiting {
			return
		}
		doneWithAction()
		if u.Bullets == 0 {
			u.Bullets = e.a.BulletMax
		}

	case pb.Action_logout:
		if u.fsm != waiting {
			return
		}
		doneWithAction()

	case pb.Action_shield:
		if u.fsm != waiting {
			return
		}
		doneWithAction()
		// No shield or shield cooling down?
		if u.NumShield == 0 || u.shieldExpiry.After(cmn.GameTime) {
			return
		}
		// Set state only, set timeout after eval resp
		u.NumShield -= 1
		u.ShieldHealth = e.a.ShieldHpMax
		u.ShieldTime = e.a.ShieldTime.Seconds()

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
			inflict(v, e.a.BulletDmg, e.a)
		} else {
			// Set timeout
			u.fsm = shotBullet
			u.shootTimeout.Reset(e.a.ShootErr)
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

		inflict(u, e.a.BulletDmg, e.a)
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

func (e *Engine) handleFov(rsp *pb.InFovResp) {
	log.Printf("Eng|Handling %v, rnd=%v\n", rsp, rsp.Rnd)
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
		inflict(me, e.a.GrenadeDmg, e.a)
	}
}

func (e *Engine) GetPlayers(i uint32) (*PlayerImpl, *PlayerImpl) {
	switch i {
	case 1:
		return &e.state[0], &e.state[1]
	case 2:
		return &e.state[1], &e.state[0]
	default:
		log.Fatal("Eng|Unknown player", i)
	}
	return nil, nil
}

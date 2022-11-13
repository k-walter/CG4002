package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

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
	rtt := (1-e.a.LPF)*e.rtt.Seconds() + e.a.LPF*en.Sub(st).Seconds()
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
	case pb.Action_none:
	case pb.Action_reload:
	case pb.Action_shoot: // OPTIMIZE clear opp's shoot/shot

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
	// Else, don't reset ShieldTime because first shielding set ShieldTime = MAX
	return p.PlayerState
}

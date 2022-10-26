package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"time"
)

type eEvalResp struct {
	*pb.State
	isLogout bool
}

func (e *eEvalResp) updateEngine(engine *Engine) bool {
	now := uint64(time.Now().UnixNano())
	updatePlayer := func(p uint32, old *cmn.PlayerState, new *pb.PlayerState) {
		// Clear shoot/shot stream
		old.Shoot = make(map[uint32]struct{})
		old.Shot = make(map[uint32]struct{})

		// Update drifitng shield time
		nextShield := now + uint64(new.ShieldTime*1e9)
		//if new.ShieldTime > 0 && cmn.AbsDiff(old.ShieldExpireNs, nextShield) > engine.ShieldErrNs {
		//	old.ShieldExpireNs = nextShield
		//}

		// Clear wrong shield
		// if eval time is zero, must have sent new shield event
		// if eval time not zero, must have exactly 1 previous shield --> didnt semd new shield event
		if old.Action == pb.Action_shield && new.Action != pb.Action_shield && new.ShieldTime == 0 {
			old.ShieldExpireNs = cmn.ShieldRst
			cmn.Pub(cmn.Event2Viz, &pb.Event{
				Player: p,
				Time:   now,
				Action: pb.Action_shieldAvailable,
			})
		}

		// Didn't predict shield
		if old.Action != pb.Action_shield && new.Action == pb.Action_shield && old.ShieldExpireNs == 0 {
			old.ShieldExpireNs = nextShield
			go waitUnshield(nextShield, &pb.Event{
				Player: p,
				Time:   nextShield,
				Action: pb.Action_shieldAvailable,
			})
		}

		// Logout?
		e.isLogout = e.isLogout || (new.Action == pb.Action_logout)

		// Copy pb
		old.PlayerState = new
	}

	// Update each player
	e.isLogout = false
	updatePlayer(1, &engine.state[0], e.P1)
	updatePlayer(2, &engine.state[1], e.P2)

	// Logout?
	if e.isLogout {
		engine.running = false
	}

	return true
}

func (e *eEvalResp) alertVizEvent() *pb.Event {
	if !e.isLogout {
		return nil
	}
	return &pb.Event{
		Player: 1,
		Time:   uint64(time.Now().UnixNano()),
		Action: pb.Action_none,
	}
}

func (e *eEvalResp) updateVizState() bool {
	return true
}

func (e *eEvalResp) updateEvalState() bool {
	// This update was received from the eval
	return false
}

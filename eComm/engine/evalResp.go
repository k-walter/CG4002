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
	updatePlayer := func(old *cmn.PlayerState, new *pb.PlayerState) {
		// Update shield time if drift
		lastShield := now - uint64(new.ShieldTime)
		if cmn.AbsDiff(old.LastShieldNs, lastShield) > engine.ShieldErrNs {
			old.LastShieldNs = lastShield
		}

		// Clear shoot/shot stream
		old.Shoot = make(map[uint32]struct{})
		old.Shot = make(map[uint32]struct{})

		// TODO if inference not working :(
		// If previous was shield, need to send shieldAvailable
		// If previous not shield, but actual is shield, send shield (assuming immediate resp)

		// Logout?
		e.isLogout = e.isLogout || (new.Action == pb.Action_logout)

		// Copy pb
		old.PlayerState = new
	}

	// Update each player
	e.isLogout = false
	updatePlayer(&engine.state[0], e.P1)
	updatePlayer(&engine.state[1], e.P2)

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

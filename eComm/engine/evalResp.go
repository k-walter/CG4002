package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"time"
)

type eEvalResp struct{ *pb.State }

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

		// Copy pb
		old.PlayerState = new
	}

	updatePlayer(&engine.state[0], e.P1)
	updatePlayer(&engine.state[1], e.P2)
	return true
}

func (e *eEvalResp) alertVizEvent() *pb.Event {
	return nil
}

func (e *eEvalResp) updateVizState() bool {
	return true
}

func (e *eEvalResp) updateEvalState() bool {
	// This update was received from the eval
	return false
}

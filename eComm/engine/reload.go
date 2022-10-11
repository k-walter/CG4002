package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
)

type eReload struct {
	*pb.Event
}

func (e *eReload) updateEngine(engine *Engine) bool {
	// Reload only if empty mag
	u, _ := engine.getStates(e.Player)
	if u.Bullets > 0 {
		return false
	}

	// RULE Unlimited magazines
	u.Bullets = cmn.BulletMax
	u.Action = pb.Action_reload

	return true
}

func (e *eReload) alertVizEvent() *pb.Event {
	// Viz no action, only update display with state
	return nil
}

func (e *eReload) updateVizState() bool {
	return true
}

func (e *eReload) updateEvalState() bool {
	return true
}

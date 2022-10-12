package engine

import (
	pb "cg4002/protos"
)

type eLogout struct {
	Event *pb.Event
}

func (e *eLogout) updateEngine(engine *Engine) bool {
	// RULE shutdown for player who signalled
	// RULE shutdown ends game for both, even if 1 player signalled
	u, _ := engine.getStates(e.Event.Player)
	u.Action = pb.Action_logout

	// DO NOT turn off engine. Check with eval first
	return true
}

func (e *eLogout) alertVizEvent() *pb.Event {
	return nil
}

func (e *eLogout) updateVizState() bool {
	return false
}

func (e *eLogout) updateEvalState() bool {
	return true
}

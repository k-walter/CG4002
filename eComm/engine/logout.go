package engine

import (
	pb "cg4002/protos"
)

type eLogout struct {
	Event *pb.Event
}

func (e *eLogout) updateEngine(engine *Engine) bool {
	// RULE shutdow for both player, even if 1 signal
	u, v := engine.getStates(e.Event.Player)
	u.Action = pb.Action_logout
	v.Action = pb.Action_logout
	engine.running = false
	return true
}

func (e *eLogout) alertVizEvent() *pb.Event {
	return e.Event
}

func (e *eLogout) updateVizState() bool {
	return true
}

func (e *eLogout) updateEvalState() bool {
	return true
}

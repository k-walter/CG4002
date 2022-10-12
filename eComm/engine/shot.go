package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
)

type eShot struct {
	*pb.Event
}

func (e *eShot) updateEngine(engine *Engine) bool {
	v, u := engine.getStates(e.Player)

	// Add to shot stream and match with shoot
	cmn.EXIT_UNLESS(len(v.Shoot) == 0 || v.Shoot[len(v.Shoot)-1] <= e.Time)
	v.Shot = append(v.Shot, e.Time)
	if !matchShot(e.Player, v) {
		return false
	}

	u.Action = pb.Action_shoot
	return true
}

func (e *eShot) alertVizEvent() *pb.Event {
	return e.Event
}

func (e *eShot) updateVizState() bool {
	return true
}

func (e *eShot) updateEvalState() bool {
	return true
}

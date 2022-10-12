package engine

import (
	pb "cg4002/protos"
	"log"
)

type eShot struct {
	*pb.Event
}

func (e *eShot) updateEngine(engine *Engine) bool {
	v, u := engine.getStates(e.Player)

	// Add to shoot stream and match with shot
	if _, fnd := v.Shot[e.ShootID]; fnd {
		log.Fatalf("should not have duplicate shootID %v\n", e.ShootID)
	}
	v.Shot[e.ShootID] = struct{}{}
	if !matchShot(e.ShootID, e.Player, v) {
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

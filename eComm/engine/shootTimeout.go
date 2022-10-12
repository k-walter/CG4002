package engine

import (
	pb "cg4002/protos"
)

type eShootTimeout struct {
	*pb.Event
}

func (s *eShootTimeout) updateEngine(e *Engine) bool {
	// If cleared, do nothing
	u, v := e.getStates(s.Player)
	if _, fnd := v.Shoot[s.ShootID]; !fnd {
		return false
	}

	// Clear shoot
	delete(v.Shoot, s.ShootID)

	u.Action = pb.Action_shoot
	return true
}

func (s *eShootTimeout) alertVizEvent() *pb.Event {
	// Already sent event in shoot
	return nil
}

func (s *eShootTimeout) updateVizState() bool {
	return true
}

func (s *eShootTimeout) updateEvalState() bool {
	return true
}

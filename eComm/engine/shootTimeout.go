package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
)

type eShootTimeout struct {
	*pb.Event
}

func (s *eShootTimeout) updateEngine(e *Engine) bool {
	// Find for uncleared shoot
	u, v := e.getStates(s.Player)
	idx := cmn.BinarySearch(v.Shoot, s.Time)

	// If cleared, do nothing
	if idx != -1 {
		return false
	}

	// Clear preceding shoots
	v.Shoot = v.Shoot[idx+1:]

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

package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
)

type eGrenade struct {
	*pb.Event
}

func (e *eGrenade) updateEngine(engine *Engine) bool {
	// Have grenades?
	u, _ := engine.getStates(e.Player)
	if u.Grenades == 0 {
		return false
	}
	u.Grenades--

	// No timeout needed. Why?
	// 1. grenade countdown on viz
	// 2. checking inFov should return true/false regardless, we update eval then
	u.Action = pb.Action_grenade
	return true
}

func (e *eGrenade) alertVizEvent() *pb.Event {
	// RULE check if hit/miss
	return e.Event
}

func (e *eGrenade) updateVizState() bool {
	// Either is fine, might have 2 updates if grenaded
	return false
}

func (e *eGrenade) updateEvalState() bool {
	// RULE check hit/miss < update
	return false
}

type EGrenaded struct {
	*pb.InFovResp
}

func (e *EGrenaded) updateEngine(engine *Engine) bool {
	v, u := engine.getStates(e.Player)
	u.Action = pb.Action_grenade

	// Hit
	if e.InFov {
		inflict(e.Player, v, cmn.GrenadeDmg)
	}
	return true
}

func (e *EGrenaded) alertVizEvent() *pb.Event {
	// Visualizer knows if hit/miss
	return nil
}

func (e *EGrenaded) updateVizState() bool {
	return true
}

func (e *EGrenaded) updateEvalState() bool {
	// RULE only now do we know if hit/miss
	return true
}

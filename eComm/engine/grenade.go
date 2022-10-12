package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
)

type eGrenade struct {
	*pb.Event
	didGrenade bool
}

func (e *eGrenade) updateEngine(engine *Engine) bool {
	u, _ := engine.getStates(e.Player)
	u.Action = pb.Action_grenade

	// Have grenades?
	e.didGrenade = u.Grenades != 0
	if e.didGrenade {
		u.Grenades--
	}

	// No timeout needed. Why?
	// 1. grenade countdown on viz
	// 2. checking inFov should return true/false regardless, we update eval then
	return true
}

func (e *eGrenade) alertVizEvent() *pb.Event {
	// Send only if didGrenade
	if !e.didGrenade {
		return nil
	}

	// RULE check if hit/miss
	return e.Event
}

func (e *eGrenade) updateVizState() bool {
	// Either is fine, might have 2 updates if grenaded
	return !e.didGrenade
}

func (e *eGrenade) updateEvalState() bool {
	// RULE check hit/miss < update
	// If no grenade left, just update now
	// Else, wait for viz event to get back
	return !e.didGrenade
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

package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"time"
)

type eShield struct {
	*pb.Event
}

func (e *eShield) updateEngine(engine *Engine) bool {
	// RULE No shields left in this life?
	u, _ := engine.getStates(e.Player)
	if u.NumShield == 0 {
		return false
	}

	// RULE Cooldown (across lifetimes)?
	if u.LastShieldNs+cmn.ShieldNs > e.Time {
		return false
	}

	// Shield up
	u.Action = pb.Action_shield
	u.ShieldHealth = cmn.ShieldHpMax
	u.LastShieldNs = e.Time
	u.NumShield--

	// Shield timeout
	go func() {
		end := time.Unix(0, int64(e.Time+cmn.ShieldNs))
		time.Sleep(time.Until(end))
		cmn.Pub(cmn.Unshield2Eng, &eUnshield{e.Event})
	}()
	return true
}

func (e *eShield) alertVizEvent() *pb.Event {
	return e.Event
}

func (e *eShield) updateVizState() bool {
	return true
}

func (e *eShield) updateEvalState() bool {
	return true
}

// Called when timeout / revived
type eUnshield struct{ *pb.Event }

func (e *eUnshield) updateEngine(engine *Engine) bool {
	// Already unshielded?
	u, _ := engine.getStates(e.Player)
	if u.LastShieldNs != e.Time {
		return false
	}

	e.Action = pb.Action_shieldAvailable
	u.ShieldHealth = 0
	return true
}

func (e *eUnshield) alertVizEvent() *pb.Event {
	return e.Event
}

func (e *eUnshield) updateVizState() bool {
	return true
}

func (e *eUnshield) updateEvalState() bool {
	// RULE no need to alert
	return false
}

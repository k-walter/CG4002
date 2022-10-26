package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"log"
	"time"
)

type eShield struct {
	*pb.Event
}

func (e *eShield) updateEngine(engine *Engine) bool {
	u, _ := engine.getStates(e.Player)
	u.Action = pb.Action_shield

	// RULE No shields left in this life?
	if u.NumShield == 0 {
		return true
	}

	// RULE Cooldown (across lifetimes)?
	if u.ShieldExpireNs > e.Time {
		return true
	}

	// Shield up
	u.ShieldHealth = cmn.ShieldHpMax
	u.ShieldExpireNs = uint64(time.Now().UnixNano()) + cmn.ShieldNs // RULE start from when engine detects
	u.NumShield--

	// Shield timeout
	go waitUnshield(u.ShieldExpireNs, &pb.Event{
		Player: e.Player,
		Time:   u.ShieldExpireNs,
		Action: pb.Action_shieldAvailable,
	})

	return true
}

func waitUnshield(end uint64, event *pb.Event) {
	t := time.Unix(0, int64(end))
	time.Sleep(time.Until(t))
	log.Println("Trying to unshield")
	cmn.Pub(cmn.Event2Eng, event)
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
	// Already unshielded by reviving?
	// NOTE if evalResp drifts ShieldExpireNs, can't compare ==
	u, _ := engine.getStates(e.Player)
	if u.ShieldExpireNs != e.Time {
		return false
	}

	e.Action = pb.Action_shieldAvailable
	u.ShieldHealth = 0
	u.ShieldExpireNs = cmn.ShieldRst
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

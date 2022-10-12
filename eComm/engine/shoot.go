package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"log"
	"time"
)

type eShoot struct {
	*pb.Event
	matchedShot bool
}

func (e *eShoot) updateEngine(engine *Engine) bool {
	u, v := engine.getStates(e.Player)
	if u.Bullets == 0 {
		return false
	}
	u.Bullets -= 1
	u.Action = pb.Action_shoot

	// Add to shoot stream and match with shot
	if _, fnd := v.Shoot[e.ShootID]; fnd {
		log.Fatalf("should not have duplicate shootID %v\n", e.ShootID)
	}
	v.Shoot[e.ShootID] = struct{}{}
	e.matchedShot = matchShot(e.ShootID, 0b11^e.Player, v)

	// Set miss timeout to update state after
	if !e.matchedShot {
		go e.waitForShot(engine.chShoot)
	} else {
		// Shot, not idiomatic ><
		cmn.Pub(cmn.Event2Viz, &pb.Event{
			Player: 0b11 ^ e.Player,
			Time:   e.Time,
			Action: pb.Action_shot,
		})
	}

	return true
}

func (e *eShoot) alertVizEvent() *pb.Event {
	// Shoot event
	return e.Event
}

func (e *eShoot) updateVizState() bool {
	return e.matchedShot
}

func (e *eShoot) updateEvalState() bool {
	return e.matchedShot
}

func (e *eShoot) waitForShot(ch chan *eShootTimeout) {
	end := time.Unix(0, int64(e.Time+cmn.ShootErrNs))
	time.Sleep(time.Until(end))
	ch <- &eShootTimeout{Event: e.Event}
}

// Match then clear preceding shoot/shot, else no updates
// Goal: update eval when match (should clear preceding shoots), or when timeout (shootTimeout clears preceding shoots)
func matchShot(shootID uint32, u uint32, s *cmn.PlayerState) bool {
	// Match shoot and shot
	if _, found := s.Shoot[shootID]; !found {
		log.Printf("could not find shooter, %v\n", s.Shoot)
		return false
	}
	if _, found := s.Shot[shootID]; !found {
		log.Printf("could not find victim, %v\n", s.Shot)
		return false
	}

	// inflict dmg
	inflict(u, s, cmn.BulletDmg)

	return true
}

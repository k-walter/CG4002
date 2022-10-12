package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
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

	// Add to shoot stream and match with shot
	cmn.EXIT_UNLESS(len(v.Shoot) == 0 || v.Shoot[len(v.Shoot)-1] <= e.Time)
	v.Shoot = append(v.Shoot, e.Time)
	u.Action = pb.Action_shoot

	// Shot arrived?
	e.matchedShot = matchShot(0b11^e.Player, v)

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
func matchShot(u uint32, s *cmn.PlayerState) bool {
	isValid := func(i, j int) bool {
		return 0 <= i && i < len(s.Shoot) && 0 <= j && j < len(s.Shot)
	}
	isMatch := func(i, j int) bool {
		return cmn.AbsDiff(s.Shoot[i], s.Shot[j]) <= cmn.ShootErrNs
	}
	isAfter := func(j, i int) bool {
		return s.Shoot[i]+cmn.ShootErrNs < s.Shot[j]
	}

	// Match forwards
	i, j := 0, 0
	for ; isValid(i, j); i++ {
		for ; isValid(i, j) && !isMatch(i, j) && !isAfter(j, i); j++ {
		}
		if isValid(i, j) && isMatch(i, j) {
			break
		}
	}
	if !(isValid(i, j) && isMatch(i, j)) {
		return false
	}

	// Clear match and preceding
	s.Shoot = s.Shoot[i+1:]
	s.Shot = s.Shot[j+1:]

	// inflict dmg
	inflict(u, s, cmn.BulletDmg)

	return true
}

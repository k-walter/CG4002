package viz

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"log"
	"time"
)

func (viz *Visualizer) diffPlayer(i int, s *pb.State, t *cmn.EvalResp) (evs []*pb.Event) {
	u, v := getPlayer(s, i), getPlayer(t.State, i)
	appendAction := func() {
		evs = append(evs, &pb.Event{Player: uint32(i), Action: v.Action})
	}

	// Revive and refresh shield timeout
	shTimer := viz.shieldTimeout[i-1]
	if u.NumDeaths != v.NumDeaths && shTimer.Stop() {
		evs = append(evs, &pb.Event{Player: uint32(i), Action: pb.Action_shieldAvailable})
		cmn.Drain(shTimer.C)
	}

	switch v.Action {
	case pb.Action_grenade:
		// Check threw
		aliveAndThrew := u.NumDeaths == v.NumDeaths && u.Grenades > v.Grenades
		reviveAndThrew := u.NumDeaths != v.NumDeaths && v.Grenades < cmn.GrenadeMax
		threw := aliveAndThrew || reviveAndThrew
		if !threw {
			return
		}
		appendAction()

		// Check missed
		if !missed(0b11^i, s, t.State) {
			evs = append(evs, &pb.Event{Player: uint32(0b11 ^ i), Action: pb.Action_grenaded})
		}

	case pb.Action_shoot:
		// Check shoot
		aliveAndShoot := u.NumDeaths == v.NumDeaths && u.Bullets > v.Bullets
		reviveAndShoot := u.NumDeaths != v.NumDeaths && v.Bullets < cmn.BulletMax
		shoot := aliveAndShoot || reviveAndShoot
		if !shoot {
			return
		}
		appendAction()

		// Check missed
		if !missed(0b11^i, s, t.State) {
			evs = append(evs, &pb.Event{Player: uint32(0b11 ^ i), Action: pb.Action_shot})
		}

	case pb.Action_reload:
		if u.Bullets == 0 && v.Bullets == cmn.BulletMax {
			appendAction()
		}

	case pb.Action_shield:
		aliveAndShield := u.NumDeaths == v.NumDeaths && u.NumShield > v.NumShield
		reviveAndShield := u.NumDeaths != v.NumDeaths && v.NumShield < cmn.ShieldMax
		shield := aliveAndShield || reviveAndShield
		if !shield {
			return
		}
		appendAction()
		shTimer.Reset(t.Time.Add(cmn.ShieldTime).Sub(time.Now()))

	case pb.Action_logout:
		evs = append(evs, &pb.Event{
			Player: uint32(i),
			Action: v.Action,
		})

	case pb.Action_none:
		return

	case pb.Action_shieldAvailable:
		fallthrough
	case pb.Action_shot:
		fallthrough
	case pb.Action_grenaded:
		fallthrough
	case pb.Action_checkFov:
		log.Fatal("Invalid player state", v.Action)

	default:
		log.Fatal("Unhandled action")
	}

	return
}

func missed(i int, s *pb.State, t *pb.State) bool {
	u, v := getPlayer(s, i), getPlayer(t, i)
	if u.NumDeaths != v.NumDeaths {
		cmn.EXIT_UNLESS(u.NumDeaths < v.NumDeaths)
		return false
	}
	return u.Hp == v.Hp
}

func getPlayer(s *pb.State, i int) *pb.PlayerState {
	switch i {
	case 1:
		return s.P1
	case 2:
		return s.P2
	default:
		log.Fatal("Unknown player ", i)
	}
	return nil
}

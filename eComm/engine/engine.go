package engine

import (
	cmn "cg4002/eComm/common"
	pb "cg4002/protos"
	"log"
	"math"
	"time"
)

type State = [cmn.NPlayer]cmn.PlayerState
type Engine struct {
	state   State
	chEvent chan *pb.Event
	chEval  chan *pb.State
}

func Make(*cmn.Arg) *Engine {
	e := Engine{
		state:   State{cmn.NewState(), cmn.NewState()},
		chEvent: make(chan *pb.Event, cmn.ChSz),
		chEval:  make(chan *pb.State, cmn.ChSz),
	}

	// Subscribe to channels
	cmn.Sub(cmn.Event2Eng, func(i interface{}) {
		// TODO pq buffer, persistence
		go func(i *pb.Event) { e.chEvent <- i }(i.(*pb.Event))
	})
	cmn.Sub(cmn.State2Eng, func(i interface{}) {
		go func(i *pb.State) { e.chEval <- i }(i.(*pb.State))
	})

	return &e
}

func (e *Engine) Run() {
	for {
		// Serialise channels
		// TODO add timestamp
		select {
		case event := <-e.chEvent:
			log.Println("engine|Received event", event.Action.String())
			if stateChanged := e.handleEvent(event); stateChanged {
				// TODO copy new state
				state := &e.state
				cmn.Pub(cmn.State2Eval, state)
				cmn.Pub(cmn.State2Viz, state) // TODO remove if often wrong/jerky
			}

		case state := <-e.chEval:
			log.Println("engine|Update with eval's truth")
			e.updateState(state)
			// TODO copy new state
			cmn.Pub(cmn.State2Viz, state)
		}
	}
}

func (e *Engine) updateState(s *pb.State) {
	now := float64(time.Now().UnixNano())
	updatePlayer := func(old *cmn.PlayerState, new *pb.PlayerState) {
		// Update shield time if drift
		lastShield := now - new.ShieldTime
		if math.Abs(float64(old.LastShieldNs)-lastShield) > cmn.ShieldErrNs {
			old.LastShieldNs = uint64(lastShield)
		}

		// Clear shoot/shot stream
		old.Shoot = old.Shoot[:0]
		old.Shot = old.Shot[:0]

		// Copy pb
		old.PlayerState = new
	}

	updatePlayer(&e.state[0], s.P1)
	updatePlayer(&e.state[1], s.P2)
}

func (e *Engine) Close() {
	// TODO Close channels
}

func (e *Engine) handleEvent(event *pb.Event) bool {
	// TODO should update eval grenade/shoot first?
	// TODO refactor handling of personal state, visualizer, eval into strategies

	// Get <from, to>
	u, v := func() (*cmn.PlayerState, *cmn.PlayerState) {
		switch event.Player {
		case 1:
			return &e.state[0], &e.state[1]
		case 2:
			return &e.state[1], &e.state[0]
		default:
			log.Fatal("Unknown player ", event.Player)
		}
		return nil, nil
	}()

	// Decide action
	switch event.Action {
	case pb.Action_shoot:
		if u.Bullets == 0 {
			return false
		}
		u.Bullets -= 1

		// Add to shoot stream and match with shot
		v.Shoot = append(v.Shoot, event.Time)

		// RULE Viz to play sound, even if not shot
		cmn.Pub(cmn.Event2Viz, event)

		// inflict dmg
		if doShoot(v) {
			inflict(v.PlayerState, cmn.BulletDmg)
		}
		return true

	case pb.Action_shot:
		// Add to shot stream and match with shoot
		u.Shot = append(u.Shot, event.Time)
		if !doShoot(u) {
			return false
		}

		// inflict dmg
		inflict(u.PlayerState, cmn.BulletDmg)
		// RULE Viz to show damaged
		cmn.Pub(cmn.Event2Viz, event)
		return true

	case pb.Action_grenade:
		if u.Grenades == 0 {
			return false
		}
		u.Grenades -= 1

		// RULE Check opponent grenaded
		// RULE viz to play explosion, even if not grenaded
		cmn.Pub(cmn.Event2Viz, event)

		// Wait until grenaded
		// TODO what if not grenade? want response regardless
		return true

	case pb.Action_grenaded:
		inflict(u.PlayerState, cmn.GrenadeDmg)
		// RULE viz to show damage
		cmn.Pub(cmn.Event2Viz, event)
		return true

	case pb.Action_reload:
		// Reload only if empty mag
		if u.Bullets > 0 {
			return false
		}

		// RULE Unlimited magazines
		u.Bullets = cmn.BulletMax

		// RULE viz no action, only update display

	case pb.Action_shield:
		// RULE No shields left in this life?
		if u.NumShield == 0 {
			return false
		}

		// RULE Cooldown (across lifetimes)?
		if u.LastShieldNs+cmn.ShieldNs > event.Time {
			return false
		}

		// Shield up
		u.ShieldHealth = cmn.ShieldHpMax
		u.LastShieldNs = event.Time
		u.NumShield -= 1

		// RULE Viz to show shield
		cmn.Pub(cmn.Event2Viz, event)
		return true

	case pb.Action_logout:
		// TODO Viz to show logout?
		cmn.Pub(cmn.Event2Viz, event)
		return true

	case pb.Action_none:
		return false

	default:
		log.Fatal("Unknown action", event.Action)
	}
	return false
}

func doShoot(s *cmn.PlayerState) bool {
	isEmpty := func() bool {
		return len(s.Shoot) == 0 || len(s.Shot) == 0
	}

	// Dequeue non-matching shoot/shot
	// Case shoot << shot: shoot but missed
	// Case shot << shoot: tried to shoot but no bullets (or lost axn)
	for !isEmpty() && math.Abs(float64(s.Shoot[0])-float64(s.Shot[0])) > cmn.ShootErrNs {
		if s.Shoot[0] < s.Shot[0] {
			s.Shoot = s.Shoot[1:]
		} else {
			s.Shot = s.Shot[1:]
		}
	}
	if isEmpty() {
		return false
	}

	// Match 1st shoot & shot
	s.Shoot = s.Shoot[1:]
	s.Shot = s.Shot[1:]

	return true
}

func inflict(player *pb.PlayerState, dmg uint32) {
	if player.Hp+player.ShieldHealth <= dmg {
		// Died
		player.NumDeaths += 1
		// TODO send died?
		// TODO auto revive
	} else if player.ShieldHealth <= dmg {
		// Dmg shield+health
		player.Hp -= dmg - player.ShieldHealth
		player.ShieldHealth = 0
	} else {
		// Dmg shield
		player.ShieldHealth -= dmg
	}
}

package engine

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"log"
)

type Engine struct {
	state   *pb.State
	chEvent chan *pb.Event
	chEval  chan *pb.State
}

func Make(*common.Arg) *Engine {
	e := Engine{
		state: &pb.State{
			P1: &pb.PlayerState{},
			P2: &pb.PlayerState{},
		},
		chEvent: make(chan *pb.Event, common.ChSz),
		chEval:  make(chan *pb.State, common.ChSz),
	}

	// Subscribe to channels
	common.Sub(common.Event2Eng, func(i interface{}) {
		go func(i *pb.Event) { e.chEvent <- i }(i.(*pb.Event))
	})
	common.Sub(common.State2Eng, func(i interface{}) {
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
			e.handleEvent(event)
			common.Pub(common.State2Eval, e.state)
			common.Pub(common.State2Viz, e.state) // TODO remove if often wrong/jerky

		case state := <-e.chEval:
			log.Println("engine|Update with eval's truth")
			e.state = state
			common.Pub(common.State2Viz, e.state)
		}
	}
}

func (e *Engine) Close() {
	// TODO Close channels
}

func (e *Engine) handleEvent(event *pb.Event) {
	// Get player
	player := func() *pb.PlayerState {
		switch event.Player {
		case 1:
			return e.state.P1
		case 2:
			return e.state.P2
		default:
			log.Fatal("Unknown player ", event.Player)
		}
		return nil
	}()

	// Decide action
	// TODO strategy pattern
	switch event.Action {
	case pb.Action_shoot:
		if player.Bullets == 0 {
			return
		}
		player.Bullets -= 1
		// TODO add to shoot stream
		// TODO resolve shoot-shot streams, inflict dmg, then fwd to viz

	case pb.Action_shot:
		const dmg = 10 // TODO clarify
		// NOTE Always corresponds to some shoot
		// TODO add to shot stream
		// TODO resolve shoot-shot streams, inflict dmg, then fwd to viz

	case pb.Action_grenade:
		if player.Grenades == 0 {
			return
		}
		player.Grenades -= 1

		// Check opponent grenaded
		// TODO viz to take req as action regardless
		common.Pub(common.InFovReq, &pb.InFovMessage{
			Player: 0b11 ^ event.Player, // other player
			Time:   event.Time,
			InFov:  true,
		})

	case pb.Action_grenaded:
		const dmg = 10 // TODO clarify
		inflict(player, dmg)

	case pb.Action_reload:
		fallthrough

	case pb.Action_shield:
		// TODO check timeout
		fallthrough

	case pb.Action_logout:
		// TODO send logout signals
		fallthrough

	case pb.Action_none:
		return

	default:
		log.Fatal("Unknown action", event.Action)
	}
}

func inflict(player *pb.PlayerState, dmg uint32) {
	if player.Hp+player.ShieldHealth <= dmg {
		// Died
		player.NumDeaths += 1
		// TODO send died?
	} else if player.ShieldHealth <= dmg {
		// Dmg shield+health
		player.Hp -= dmg - player.ShieldHealth
		player.ShieldHealth = 0
	} else {
		// Dmg shield
		player.ShieldHealth -= dmg
	}
}

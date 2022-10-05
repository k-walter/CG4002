package engine

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"fmt"
	"google.golang.org/grpc"
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


package engine

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Engine struct {
	state pb.State
	chRelay chan struct{}
	chEval chan *pb.State
}

func Make(a *common.Arg) Engine {
	// Subscribe to channels
	RelayToEngine Topic = iota
	EvalToEngine

    return Engine{state: pb.State{
		P1: &pb.PlayerState{},
		P2: &pb.PlayerState{},
	}}
}

func (s *RelayServer) Run() {
	// Multiplex channels
}

func (s *RelayServer) Close() {
	// TODO Close channels
}


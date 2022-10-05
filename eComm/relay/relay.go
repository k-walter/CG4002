package relay

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

type RelayServer struct {
	// From Relay
	pb.UnimplementedRelayServer
	lis net.Listener
}

func Make(a *common.Arg) *RelayServer {
	// From relay
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.RelayPort))
	if err != nil {
		log.Fatal(err)
	}

	// To pynq
	pyConn, err := grpc.Dial(fmt.Sprintf(":%v", a.ToPynqPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	return &RelayServer{
		lis:    lis,
		py:     pb.NewToPynqClient(pyConn),
		pyConn: pyConn,
	}
}

func (s *RelayServer) Run() {
	g := grpc.NewServer()
	pb.RegisterFromRelayServer(g, s)

	if err := g.Serve(s.lis); err != nil {
		log.Fatal(err)
	}
}

func (s *RelayServer) Close() {
	_ = s.lis.Close()
	_ = s.pyConn.Close()
}

func (s *RelayServer) Gesture(c context.Context, d *pb.SensorData) (*emptypb.Empty, error) {
	log.Println("relay|Received gesture")
	d.Time = uint32(time.Now().Nanosecond())
	common.Pub(common.Data2Pynq, d)
	return &emptypb.Empty{}, nil
}

func (s *RelayServer) Shoot(c context.Context, e *pb.Event) (*emptypb.Empty, error) {
	log.Println("relay|Received shoot")
	common.Pub(common.EventToEngine, pb.Event{
		Player: d.Player,
		Time:   uint32(time.Now().Nanosecond()),
		Action: pb.Action_shoot,
	})
	return &emptypb.Empty{}, nil
}

func (s *RelayServer) Shot(c context.Context, e *pb.Event) (*emptypb.Empty, error) {
	log.Println("relay|Received shot")
	common.Pub(common.EventToEngine, pb.Event{
		Player: d.Player,
		Time:   uint32(time.Now().Nanosecond()),
		Action: pb.Action_shot,
	})
	return &emptypb.Empty{}, nil
}

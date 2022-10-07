package relay

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

type Server struct {
	// From Relay
	pb.UnimplementedRelayServer
	lis net.Listener
}

func Make(a *common.Arg) *Server {
	// From relay
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.RelayPort))
	if err != nil {
		log.Fatal(err)
	}

	return &Server{lis: lis}
}

func (s *Server) Run() {
	g := grpc.NewServer()
	pb.RegisterRelayServer(g, s)
	if err := g.Serve(s.lis); err != nil {
		log.Fatal(err)
	}
	log.Println("running relay")

}

func (s *Server) Close() {
	_ = s.lis.Close()
}

func (s *Server) Gesture(c context.Context, d *pb.SensorData) (*emptypb.Empty, error) {
	log.Println("relay|Received gesture")
	d.Time = uint64(time.Now().UnixNano())
	common.Pub(common.Data2Pynq, d)
	return &emptypb.Empty{}, nil
}

func (s *Server) Shoot(c context.Context, e *pb.Event) (*emptypb.Empty, error) {
	log.Println("relay|Received shoot")

	// Verification
	if !(1 <= e.Player && e.Player <= 2) {
		return nil, status.Error(codes.Unknown, "Player must be 1/2")
	}
	if e.Action != pb.Action_shoot {
		return nil, status.Error(codes.Unknown, "Shoot() called with non-shoot action")
	}

	// Forward to engine
	e.Time = uint64(time.Now().UnixNano())
	common.Pub(common.Event2Eng, e)

	return &emptypb.Empty{}, nil
}

func (s *Server) Shot(c context.Context, e *pb.Event) (*emptypb.Empty, error) {
	log.Println("relay|Received shot")

	// Verification
	if !(e.Player == 1 || e.Player == 2) {
		return nil, status.Error(codes.Unknown, "Player must be 1/2")
	}
	if e.Action != pb.Action_shot {
		return nil, status.Error(codes.Unknown, "Shot() called with non-shot action")
	}

	// Forward to engine
	e.Time = uint64(time.Now().UnixNano())
	common.Pub(common.Event2Eng, e)

	return &emptypb.Empty{}, nil
}

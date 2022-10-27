package relay

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
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
	// From engine
	chRnd chan uint32
	// Metrics
	nextMetric time.Time
	hz         int
}

func Make(a *common.Arg) *Server {
	s := Server{
		lis:        nil,
		chRnd:      make(chan uint32, common.ChSz),
		nextMetric: time.Now(),
		hz:         0,
	}

	// From relay
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.RelayPort))
	if err != nil {
		log.Fatal(err)
	}
	s.lis = lis

	return &s
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

func (s *Server) GetRound(_ *emptypb.Empty, stream pb.Relay_GetRoundServer) error {
	for rnd := range s.chRnd {
		err := stream.Send(&pb.RndResp{
			Rnd: rnd,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (s *Server) Gesture(stream pb.Relay_GestureServer) error {
	defer stream.SendAndClose(&emptypb.Empty{})
	for {
		d, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		// measure data rate using tumbling window
		s.hz += 1
		if now := time.Now(); !s.nextMetric.Before(now) {
			log.Println("Relay|Gesture rate = ", s.hz)
			s.hz = 0
			s.nextMetric = now.Add(time.Second)
		}

		// Send each data separately
		// OPTIMIZE stream pynq
		d.Time = uint64(time.Now().UnixNano())
		common.PubFull(common.Data2Pynq, d, false)
	}

	log.Println("relay|Closed gesture")
	return nil
}

func (s *Server) Shoot(stream pb.Relay_ShootServer) error {
	defer stream.SendAndClose(&emptypb.Empty{})
	for {
		e, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("relay|Received shoot", e.ShootID)

		// Verification
		if !(1 <= e.Player && e.Player <= 2) {
			return status.Error(codes.Unknown, "Player must be 1/2")
		}
		if e.Action != pb.Action_shoot {
			return status.Error(codes.Unknown, "Shoot() called with non-shoot action")
		}

		// Forward to engine
		e.Time = uint64(time.Now().UnixNano())
		common.Pub(common.Event2Eng, e)
	}

	log.Println("relay|Closed shoot")
	return nil
}

func (s *Server) Shot(stream pb.Relay_ShotServer) error {
	defer stream.SendAndClose(&emptypb.Empty{})
	for {
		e, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("relay|Received shot")

		// Verification
		if !(1 <= e.Player && e.Player <= 2) {
			return status.Error(codes.Unknown, "Player must be 1/2")
		}
		if e.Action != pb.Action_shoot {
			return status.Error(codes.Unknown, "Shot() called with non-shot action")
		}

		// Forward to engine
		e.Time = uint64(time.Now().UnixNano())
		common.Pub(common.Event2Eng, e)
	}

	log.Println("relay|Closed shoot")
	return nil
}

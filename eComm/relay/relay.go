package relay

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	// From Relay
	pb.UnimplementedRelayServer
	lis net.Listener

	// From engine
	chRnd chan common.RoundT

	// Bookeeping
	nextMetric time.Time
}

func Make(a *common.Arg) *Server {
	s := Server{
		lis:        nil,
		chRnd:      common.Sub[common.RoundT](common.ERound),
		nextMetric: time.Now(),
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
	log.Println("Relay|running")
	defer log.Println("Relay|stopped")
	if err := g.Serve(s.lis); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Close() {
	_ = s.lis.Close()
}

func (s *Server) GetRound(_ *emptypb.Empty, stream pb.Relay_GetRoundServer) error {
	log.Println("Relay|getRound started")
	defer log.Println("Relay|Closed getRound")
	for rnd := range s.chRnd {
		err := stream.Send(&pb.RndResp{
			Rnd: uint32(rnd),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Gesture(stream pb.Relay_GestureServer) error {
	log.Println("Relay|gesture started")
	defer log.Println("Relay|Closed gesture")
	defer stream.SendAndClose(&emptypb.Empty{})
	p1, p2 := 0, 0
	for {
		d, err := stream.Recv()
		if err != nil {
			return err
		}

		// measure data rate using tumbling window
		if d.Player == 1 {
			p1 += 1
		} else {
			p2 += 1
		}
		if now := time.Now(); !s.nextMetric.After(now) {
			log.Printf("Relay|Gesture rate=%v, p1=%v, p2=%v\n", p1+p2, p1, p2)
			p1, p2 = 0, 0
			s.nextMetric = now.Add(time.Second)
		}

		// Send each data separately
		// OPTIMIZE stream pynq
		d.Time = common.TimeToNs(time.Now())
		common.Pub(common.EData, d)
	}

	return nil
}

func (s *Server) Shoot(stream pb.Relay_ShootServer) error {
	log.Println("Relay|shoot started")
	defer log.Println("Relay|Closed shoot")
	defer stream.SendAndClose(&emptypb.Empty{})
	for {
		e, err := stream.Recv()
		if err != nil {
			return err
		}

		// Verification
		if !(1 <= e.Player && e.Player <= 2) {
			return status.Error(codes.Unknown, "Player must be 1/2")
		}
		if e.Action != pb.Action_shoot {
			return status.Error(codes.Unknown, "Shoot() called with non-shoot action")
		}

		// Forward to engine
		e.Time = common.TimeToNs(time.Now())
		common.Pub(common.EEvent, e)
	}

	return nil
}

func (s *Server) Shot(stream pb.Relay_ShotServer) error {
	log.Println("Relay|shot started")
	defer log.Println("Relay|Closed shot")
	defer stream.SendAndClose(&emptypb.Empty{})
	for {
		e, err := stream.Recv()
		if err != nil {
			return err
		}

		// Verification
		if !(1 <= e.Player && e.Player <= 2) {
			return status.Error(codes.Unknown, "Player must be 1/2")
		}
		if e.Action != pb.Action_shot {
			return status.Error(codes.Unknown, "Shot() called with non-shot action")
		}

		// Forward to engine
		e.Time = common.TimeToNs(time.Now())
		common.Pub(common.EEvent, e)
	}

	return nil
}

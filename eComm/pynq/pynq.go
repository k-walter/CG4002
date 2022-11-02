package pynq

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type Client struct {
	chData chan *pb.Data
	py     pb.PynqClient
	pyConn *grpc.ClientConn
}

func Make(a *common.Arg) *Client {
	// To pynq
	pyConn, err := grpc.Dial(fmt.Sprintf(":%v", a.PynqPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	// Init client
	c := Client{
		chData: common.SubCap[*pb.Data](common.EData, 2*3*50),
		py:     pb.NewPynqClient(pyConn),
		pyConn: pyConn,
	}

	return &c
}

func (c *Client) Run() {
	// OPTIMIZE measure latency of blocking inference, then add polling/interrupt mechanism accordingly
	//t := time.NewTicker(pollIntervalMs * time.Millisecond)
	//defer t.Stop()

	for {
		// Synchronous, since sharing grpc stream
		select {
		// Send to pynq
		case sensorData := <-c.chData:
			// Forward synchronously
			if event := c.forward(sensorData); event != nil {
				common.Pub(common.EEvent, event)
			}

			// Poll pynq
			// OPTIMIZE combine with forward?
			//case <-t.C:
			//	if event := c.pollEvent(); event != nil {
			//	}
		}
	}
}

func (c *Client) Close() {
	_ = c.pyConn.Close()
}

func (c *Client) forward(d *pb.Data) *pb.Event {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ev, err := c.py.Emit(ctx, d)
	if err != nil {
		log.Fatal(err)
	}

	// Infered anything?
	if ev.Action == pb.Action_none {
		return nil
	}

	return ev
}

func (c *Client) pollEvent() *pb.Event {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	event, err := c.py.Poll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	if event == nil || event.Action == pb.Action_none {
		return nil
	}
	return event
}

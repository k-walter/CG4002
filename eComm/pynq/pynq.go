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
	"sync/atomic"
	"time"
)

type Client struct {
	chData chan *pb.Data
	py     pb.PynqClient
	pyConn *grpc.ClientConn

	// Metrics
	qSz atomic.Int32
}

const (
	pollIntervalMs = 10
)

func Make(a *common.Arg) *Client {
	// To pynq
	pyConn, err := grpc.Dial(fmt.Sprintf(":%v", a.PynqPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	// Init client
	c := Client{
		chData: make(chan *pb.Data, common.ChSz),
		py:     pb.NewPynqClient(pyConn),
		pyConn: pyConn,
		qSz:    atomic.Int32{},
	}
	c.qSz.Store(0)

	// Subscribe to data
	common.Sub(common.Data2Pynq, func(i interface{}) {
		go func(d *pb.Data) {
			// Warn if buffer growing
			c.qSz.Add(1)
			if sz := c.qSz.Load(); sz > 50 { // 50hz
				log.Printf("Warning: relay->fpga buffer size = %v", sz)
			}

			// Push to buffer
			c.chData <- d
		}(i.(*pb.Data))
	})

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
			c.qSz.Add(-1)

			// Forward synchronously
			if event := c.forward(sensorData); event != nil {
				common.Pub(common.Event2Eng, event)
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

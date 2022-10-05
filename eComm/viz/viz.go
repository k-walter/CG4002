package viz

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const (
	broker    = "tcp://broker.hivemq.com:1883"
	timeoutMs = 1000

	inFovReqTopic  = "cg4002/b7/inFovReq"
	inFovRespTopic = "cg4002/b7/inFovResp"
	updateTopic    = "cg4002/b7/update"
)

type Visualizer struct {
	// mqtt
	clnt mqtt.Client

	chState chan *pb.State
	chInFov chan *pb.InFovMessage
}

func Make(a *common.Arg) *Visualizer {
	// Connect to mqtt broker
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("sample")
	c := mqtt.NewClient(opts)
	if t := c.Connect(); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}

	// Subscribe to state updates
	chState := make(chan *pb.State, common.ChSz)
	common.Sub(common.State2Viz, func(i interface{}) {
		go func(s *pb.State) { chState <- s }(i.(*pb.State))
	})

	// Subscribe to mqtt grenade inFov req
	chInFov := make(chan *pb.InFovMessage, common.ChSz)
	common.Sub(common.InFovReq, func(i interface{}) {
		go func(i *pb.InFovMessage) { chInFov <- i }(i.(*pb.InFovMessage))
	})

	// Subscribe to mqtt grenade inFov resp
	if t := c.Subscribe(inFovRespTopic, 0, inFovRespHandler); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}

	return &Visualizer{
		clnt:    c,
		chState: chState,
		chInFov: chInFov,
	}
}

func (v *Visualizer) Close() {
	v.clnt.Disconnect(timeoutMs)
}

func (v *Visualizer) Run() {
	for {
		select { // Async send
		case state := <-v.chState:
			go v.publishState(state)
		case msg := <-v.chInFov:
			go v.inFovReq(msg)
		}
	}
}

func (v *Visualizer) publishState(s *pb.State) {
	data, err := proto.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	// From https://www.hivemq.com/docs/hivemq/4.8/control-center/information.html#retained
	// QOS = 0 (<=1), 1 (>=1), 2 (1) semantics
	if t := v.clnt.Publish(updateTopic, 1, false, data); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}
}

func (v *Visualizer) inFovReq(msg *pb.InFovMessage) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	if t := v.clnt.Publish(inFovReqTopic, 1, false, data); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}
}

func inFovRespHandler(c mqtt.Client, m mqtt.Message) {
	data := m.Payload()
	msg := pb.InFovMessage{}
	if err := proto.Unmarshal(data, &msg); err != nil {
		log.Fatal(err)
	}

	if msg.InFov {
		common.Pub(common.Event2Eng, &pb.Event{
			Player: msg.Player,
			Time:   msg.Time,
			Action: pb.Action_grenaded,
		})
	}
}

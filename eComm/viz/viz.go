package viz

import (
	"cg4002/eComm/common"
	"cg4002/eComm/engine"
	pb "cg4002/protos"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const (
	broker    = "tcp://broker.hivemq.com:1883"
	timeoutMs = 1000

	eventTopic     = "cg4002/b7/event"
	inFovRespTopic = "cg4002/b7/inFovResp"
	updateTopic    = "cg4002/b7/state"
)

type Visualizer struct {
	// mqtt
	clnt mqtt.Client

	chState chan *pb.State
	chEvent chan *pb.Event
}

func Make(*common.Arg) *Visualizer {
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
	chEvent := make(chan *pb.Event, common.ChSz)
	common.Sub(common.Event2Viz, func(i interface{}) {
		go func(i *pb.Event) { chEvent <- i }(i.(*pb.Event))
	})

	// Subscribe to mqtt grenade inFov resp
	if t := c.Subscribe(inFovRespTopic, 0, inFovRespHandler); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}

	return &Visualizer{
		clnt:    c,
		chState: chState,
		chEvent: chEvent,
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
		case event := <-v.chEvent:
			go v.publishEvent(event)
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

func (v *Visualizer) publishEvent(e *pb.Event) {
	// NOTE Grenade event doubles as InFovRequest
	data, err := protojson.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}
	if t := v.clnt.Publish(eventTopic, 1, false, data); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}
}

func inFovRespHandler(c mqtt.Client, m mqtt.Message) {
	data := m.Payload()
	msg := pb.InFovResp{}
	if err := protojson.Unmarshal(data, &msg); err != nil {
		log.Fatal(err)
	}

	common.Pub(common.Grenade2Eng, &engine.EGrenaded{InFovResp: &msg})
}

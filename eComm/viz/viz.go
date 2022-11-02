package viz

import (
	"cg4002/eComm/common"
	"cg4002/eComm/eval"
	pb "cg4002/protos"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"sync"
	"time"
)

const (
	broker    = "tcp://broker.hivemq.com:1883"
	timeoutMs = 1000

	eventTopic     = "cg4002/b7/event"
	inFovRespTopic = "cg4002/b7/inFovResp"
	stateTopic     = "cg4002/b7/state"
)

type Visualizer struct {
	// mqtt
	clnt mqtt.Client

	chState       chan *eval.EEvalResp
	chEvent       chan *pb.Event
	state         *pb.State
	shieldTimeout [2]*time.Timer
}

func Make(*common.Arg) *Visualizer {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("sample")
	v := Visualizer{
		clnt:          mqtt.NewClient(opts),
		chState:       common.Sub[*eval.EEvalResp](common.EEvalResp),
		chEvent:       common.Sub[*pb.Event](common.EEvent),
		state:         common.NewState(),
		shieldTimeout: [2]*time.Timer{time.NewTimer(0), time.NewTimer(0)}, // must be created with NewTimer
	}

	// Connect to mqtt broker
	if t := v.clnt.Connect(); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}

	// Subscribe to mqtt grenade inFov resp
	if t := v.clnt.Subscribe(inFovRespTopic, 0, fovRespHandler); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}

	return &v
}

func (v *Visualizer) Close() {
	v.clnt.Disconnect(timeoutMs)
}

func (v *Visualizer) Run() {
	// Drain timers
	common.Drain(v.shieldTimeout[0].C)
	common.Drain(v.shieldTimeout[1].C)

	for {
		select {
		case state := <-v.chState:
			v.updateState(state) // serialize updates
		case event := <-v.chEvent:
			// TODO unrestricted mode
			go v.checkFov(event) // async send
		case <-v.shieldTimeout[0].C:
		case <-v.shieldTimeout[1].C:
			// TODO go shield avail
		}
	}
}

func (v *Visualizer) updateState(s *eval.EEvalResp) {
	// TODO unrestricted mode
	var wg sync.WaitGroup

	// Publish events
	evs := append(
		v.diffPlayer(1, v.state, s),
		v.diffPlayer(2, v.state, s)...)
	for _, e := range evs {
		wg.Add(1)
		data := common.PbToJson(e.ProtoReflect())
		go common.Do(&wg, func() { v.pub(eventTopic, data) })
	}

	// Publish state
	v.state = s.State
	wg.Add(1)
	data := common.PbToJson(v.state.ProtoReflect())
	go common.Do(&wg, func() { v.pub(stateTopic, data) })

	// Publish all before return
	wg.Done()
}

func (v *Visualizer) checkFov(e *pb.Event) {
	if e.Action != pb.Action_grenade {
		return
	}

	// Make request
	req := pb.Event{
		Player: e.Player,
		Time:   e.Time,
		Rnd:    e.Rnd,
		Action: pb.Action_checkFov, // WARNING does not imply there are enough grenades
	}
	data := common.PbToJson(req.ProtoReflect())
	v.pub(eventTopic, data)
}

func (v *Visualizer) pub(t string, data []byte) {
	// From https://www.hivemq.com/docs/hivemq/4.8/control-center/information.html#retained
	// QOS = 0 (<=1), 1 (>=1), 2 (1) semantics
	if t := v.clnt.Publish(t, 1, false, data); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}
}

func fovRespHandler(c mqtt.Client, m mqtt.Message) {
	data := m.Payload()
	msg := pb.InFovResp{}
	if err := protojson.Unmarshal(data, &msg); err != nil {
		log.Fatal(err)
	}

	common.Pub(common.EInFov, &msg)
}

package viz

import (
	"cg4002/eComm/common"
	pb "cg4002/protos"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"time"
)

const (
	timeoutMs = 1000
	//broker    = "tcp://broker.hivemq.com:1883"
	//eventTopic     = "cg4002/b7/event"
	//inFovRespTopic = "cg4002/b7/inFovResp"
	//stateTopic     = "cg4002/b7/state"
)

type Visualizer struct {
	// mqtt
	clnt mqtt.Client

	chState       chan *common.EvalResp
	chEvent       chan *pb.Event
	state         *pb.State
	shieldTimeout [2]*time.Timer
	done          [2]bool

	a *common.Arg
}

func Make(a *common.Arg) *Visualizer {
	opts := mqtt.NewClientOptions().
		AddBroker(a.Broker).
		SetClientID("sample")
	v := Visualizer{
		clnt:          mqtt.NewClient(opts),
		chState:       common.Sub[*common.EvalResp](common.EEvalResp),
		chEvent:       common.Sub[*pb.Event](common.EEvent),
		state:         common.NewState(a),
		shieldTimeout: [2]*time.Timer{time.NewTimer(0), time.NewTimer(0)}, // must be created with NewTimer
		done:          [2]bool{false, false},
		a:             a,
	}

	// Connect to mqtt broker
	if t := v.clnt.Connect(); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}

	// Subscribe to mqtt grenade inFov resp
	if t := v.clnt.Subscribe(a.InFovRespTopic, 0, fovRespHandler); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
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

	// Send initial state to clear debug
	v.pubState()
	v.pubShieldAvail(1)
	v.pubShieldAvail(2)

	for {
		select {
		case state := <-v.chState:
			v.updateState(state) // serialize updates
		case event := <-v.chEvent:
			go v.checkDone(event)
			go v.checkFov(event) // async send
		case <-v.shieldTimeout[0].C:
			go v.pubShieldAvail(1)
		case <-v.shieldTimeout[1].C:
			go v.pubShieldAvail(2)
		}
	}
}

func (v *Visualizer) updateState(s *common.EvalResp) {
	// Publish events
	evs := append(
		v.diffPlayer(1, v.state, s),
		v.diffPlayer(2, v.state, s)...)
	for _, e := range evs {
		data := common.PbToJson(e.ProtoReflect())
		v.pub(v.a.EventTopic, data)
	}

	// Publish state
	v.state = s.State
	v.pubState()

	// Reset done flag
	v.done[0] = false
	v.done[1] = false
}

func (v *Visualizer) checkFov(e *pb.Event) {
	if e.Action != pb.Action_checkFov {
		return
	}
	data := common.PbToJson(e.ProtoReflect())
	v.pub(v.a.EventTopic, data)
}

func (v *Visualizer) pub(t string, data []byte) {
	// From https://www.hivemq.com/docs/hivemq/4.8/control-center/information.html#retained
	// QOS = 0 (<=1), 1 (>=1), 2 (1) semantics
	log.Println("Viz|Pub to MQTT|", string(data))
	if t := v.clnt.Publish(t, 1, false, data); t.WaitTimeout(timeoutMs*time.Millisecond) && t.Error() != nil {
		log.Fatal(t.Error())
	}
}

func (v *Visualizer) pubShieldAvail(i int) {
	e := pb.Event{
		Player: uint32(i),
		Action: pb.Action_shieldAvailable,
	}
	data := common.PbToJson(e.ProtoReflect())
	v.pub(v.a.EventTopic, data)
}

func (v *Visualizer) checkDone(e *pb.Event) {
	if v.done[e.Player-1] {
		return
	}
	switch e.Action {
	case pb.Action_grenade:
		fallthrough
	case pb.Action_reload:
		fallthrough
	case pb.Action_shoot:
		fallthrough
	case pb.Action_logout:
		fallthrough
	case pb.Action_shield:
		v.done[e.Player-1] = true
		ev := pb.Event{
			Player: e.Player,
			Time:   e.Time,
			Rnd:    e.Rnd,
			Action: pb.Action_done,
		}
		data := common.PbToJson(ev.ProtoReflect())
		v.pub(v.a.EventTopic, data)

	case pb.Action_none:
	case pb.Action_shot:
	case pb.Action_grenaded:
	case pb.Action_shieldAvailable:
	case pb.Action_checkFov:
	}
}

func (v *Visualizer) pubState() {
	data := common.PbToJson(v.state.ProtoReflect())
	v.pub(v.a.StateTopic, data)
}

func fovRespHandler(c mqtt.Client, m mqtt.Message) {
	data := m.Payload()
	msg := pb.InFovResp{}
	if err := protojson.Unmarshal(data, &msg); err != nil {
		log.Fatal(err)
	}

	common.Pub(common.EInFov, &msg)
}

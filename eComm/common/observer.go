package common

import "log"

type Topic uint16

const (
	nTopic = 7
	ChSz   = 1
)

const (
	// From relay
	Event2Eng Topic = iota // *pb.Event
	Data2Pynq              // *pb.SensorData

	// From engine
	State2Eng // *pb.State
	State2Viz // *pb.State
	Event2Viz // *pb.Event

	// From viz
	Grenade2Eng // *EGrenaded

	// From relay
	State2Eval // *pb.State
)

func (t Topic) String() string {
	switch t {
	case Event2Eng:
		return "Event2Eng"
	case Data2Pynq:
		return "Data2Pynq"
	case State2Eval:
		return "State2Eval"
	case State2Eng:
		return "State2Eng"
	case State2Viz:
		return "State2Viz"
	case Event2Viz:
		return "Event2Viz"
	case Grenade2Eng:
		return "Grenade2Eng"
	}
	log.Fatalf("unknown enum %d\n", t)
	return "unk"
}

var observer = make([][]func(interface{}), nTopic)

func MakeObserver() {
	for i := range observer {
		observer[i] = make([]func(interface{}), 0)
	}
}

// Sub is not thread safe. To be setup at the beginning
func Sub(t Topic, f func(interface{})) {
	log.Printf("init|Subscribe to %s %d\n", t, t)
	observer[t] = append(observer[t], f)
}

// Pub is thread safe and blocking. To publish asynchronously, subscriber should pass a goroutine closure
func Pub(t Topic, i interface{}) {
	PubFull(t, i, true)
}

func PubFull(t Topic, i interface{}, debug bool) {
	if debug {
		log.Println("Publish", t)
	}
	for _, f := range observer[t] {
		f(i)
	}
}

package common

import "log"

type Topic uint16

const (
	nTopic = 5
	ChSz   = 1

	RelayToEngine Topic = iota

	EngineToEval // *pb.State
	EvalToEngine // *pb.State

	EngineToVisualizer // *pb.State
	VisualizerToEngine // ??
)

func (t Topic) String() string {
	switch t {
	case RelayToEngine:
		return "RelayToEngine"
	case EngineToEval:
		return "EngineToEval"
	case EvalToEngine:
		return "EvalToEngine"
	case EngineToVisualizer:
		return "EngineToVisualizer"
	case VisualizerToEngine:
		return "VisualizerToEngine"
	}
	log.Fatal("unknown enum", t)
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
	log.Println("init|Subscribe to", t)
	observer[t] = append(observer[t], f)
}

// Pub is thread safe and blocking. To publish asynchronously, subscriber should pass a goroutine closure
func Pub(t Topic, i interface{}) {
	for _, f := range observer[t] {
		log.Println("Publish", t)
		f(i)
	}
}

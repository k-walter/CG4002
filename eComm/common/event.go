package common

import (
	pb "cg4002/protos"
	"log"
	"time"
)

type EventE uint16

const (
	EEvent EventE = iota

	nEvents
	chSz = 1000
)

func (t EventE) String() string {
	switch t {
	case EEvent:
		return "pb.Event"
	default:
		log.Fatalf("unknown enum %d\n", t)
	}
	return "unk"
}

type EventT interface {
	*pb.Event
}

var events = make([][]interface{}, nEvents)

func Sub[T EventT](e EventE) chan T {
	ch := make(chan T, chSz)
	events[e] = append(events[e], ch)
	return ch
}

func Pub[T EventT](e EventE, v T) {
	for _, ch := range events[e] {
		for !pubOne(ch.(chan T), v) {
			log.Println(e, "channel blocked. Backpressure")
			time.Sleep(time.Millisecond)
		}
	}
}

func pubOne[T EventT](ch chan T, v T) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

package common

import (
	pb "cg4002/protos"
	"log"
	"math/rand"
	"time"
)

type EventE uint16

const (
	EEvent EventE = iota
	EData
	ERound

	nEvents
	chSz = 1000
)

func (t EventE) String() string {
	switch t {
	case EEvent:
		return "pb.Event"
	case EData:
		return "pb.Data"
	case ERound:
		return "Round(RoundT)"
	default:
		log.Fatalf("unknown enum %d\n", t)
	}
	return "unk"
}

type EventT interface {
	*pb.Event | *pb.Data |
		RoundT
}

var events = make([][]interface{}, nEvents)

func Sub[T EventT](e EventE) chan T {
	ch := make(chan T, chSz)
	events[e] = append(events[e], ch)
	return ch
}

func Pub[T EventT](e EventE, v T) {
	for _, ch := range events[e] {
		for i := 1; !pubOne(ch.(chan T), v); i++ {
			log.Printf("%s channel blocked %v times. Backpressure\n", e, i)
			oneToTenMs := time.Duration(1_000_000+rand.Intn(9_000_000)) * time.Nanosecond
			time.Sleep(oneToTenMs)
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

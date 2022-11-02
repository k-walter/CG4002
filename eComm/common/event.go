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
	EInFov
	EEvalResp

	nEvents
	chSz = 100
)

func (t EventE) String() string {
	switch t {
	case EEvent:
		return "pb.Event"
	case EData:
		return "pb.Data"
	case ERound:
		return "Round(RoundT)"
	case EInFov:
		return "pb.InFovResp"
	case EEvalResp:
		return "EvalResp"
	default:
		log.Fatalf("unknown enum %d\n", t)
	}
	return "unk"
}

type EventT interface {
	*pb.Event | *pb.Data | *pb.InFovResp |
		RoundT | *EvalResp
}

var events = make([][]interface{}, nEvents)

func Sub[T EventT](e EventE) chan T {
	return SubCap[T](e, chSz)
}

func SubCap[T EventT](e EventE, cap int) chan T {
	log.Printf("Subscribe %s\n", e)
	ch := make(chan T, cap)
	events[e] = append(events[e], ch)
	return ch
}

func Pub[T EventT](e EventE, v T) {
	if shouldLog(e) {
		log.Printf("Publish %s, val=%v\n", e, v)
	}
	for _, ch := range events[e] {
		for i := 1; !pubOne(ch.(chan T), v); i++ {
			log.Printf("%s channel blocked %v times. Backpressure\n", e, i)
			oneToTenMs := time.Duration(1_000_000+rand.Intn(9_000_000)) * time.Nanosecond
			time.Sleep(oneToTenMs)
		}
	}
}

func shouldLog(e EventE) bool {
	switch e {
	case EEvent:
		return false
	case EData:
		return false
	case ERound:
		return true
	case EInFov:
		return true
	case EEvalResp:
		return true
	default:
		log.Fatal("Unhandled event type")
	}
	return false
}

func pubOne[T EventT](ch chan T, v T) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

type EvalResp struct {
	*pb.State
	Time time.Time
}

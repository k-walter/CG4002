package common

import (
	pb "cg4002/protos"
	"github.com/go-playground/assert/v2"
	"sync"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	// Init
	ch := Sub[*pb.Event](EEvent)
	assert.NotEqual(t, ch, nil)

	// Pub
	n := chSz * 2
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			Pub(EEvent, &pb.Event{})
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Sub
	timeout := time.NewTimer(time.Second)
	m := 0
	for ok := true; ok; {
		select {
		case _, ok = <-ch:
			if ok {
				m++
			}
		case <-timeout.C:
			t.Fatalf("Timeout after getting %v values", m)
		}
	}
	assert.Equal(t, n, m)
}

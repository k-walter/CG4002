package main

import (
	"cg4002/eComm/common"
	"cg4002/eComm/engine"
	// "cg4002/eComm/eval"
	"cg4002/eComm/pynq"
	"cg4002/eComm/relay"
	"cg4002/eComm/viz"
	"sync"
)

type Proc interface {
	Run()
	Close()
}

func main() {
	// Read args, else panic
	args := common.ParseArgs()

	// Setup pub sub observer pattern
	common.MakeObserver()

	// Setup all Proc, else panic
	proc := []Proc{
		engine.Make(&args),
		relay.Make(&args),
		pynq.Make(&args),
		viz.Make(&args),
		// eval.Make(&args),
	}
	defer func() { // RAII
		for _, p := range proc {
			p.Close()
		}
	}()

	// Schedule processes
	var wg sync.WaitGroup
	for _, p := range proc {
		wg.Add(1)
		go func(p Proc) {
			p.Run()
			wg.Done()
		}(p)
	}

	// Join
	wg.Wait()
}

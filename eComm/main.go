package main

import (
	"cg4002/eComm/common"
	"cg4002/eComm/engine"
	"cg4002/eComm/pynq"
	"cg4002/eComm/relay"
	"cg4002/eComm/viz"
	"log"
	"sync"
)

type Proc interface {
	Run()
	Close()
}

func main() {
	// Log us
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) //  | log.Lshortfile

	// Read args, else panic
	args := common.ParseArgs()

	// Setup all Proc, else panic
	proc := []Proc{
		engine.Make(&args),
		relay.Make(&args),
		pynq.Make(&args),
		viz.Make(&args),
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

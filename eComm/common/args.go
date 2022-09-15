package common

import (
	"flag"
	"log"
)

type Arg struct {
	EvalPort  uint
	EvalKey   string
	RelayPort uint
	PynqPort  uint
}

func ParseArgs() Arg {
	// Bind a
	a := Arg{}
	flag.UintVar(&a.EvalPort, "evalPort", 8080, "Eval server port")
	flag.StringVar(&a.EvalKey, "evalKey", "PLSPLSPLSPLSWORK", "Symmetric key for eval server")
	flag.UintVar(&a.RelayPort, "relayPort", 8081, "Port to be SSH forwarded to relay")
	flag.UintVar(&a.PynqPort, "pynqPort", 8082, "Port to send/poll signals")
	flag.Parse()

	// Checks
	if len(a.EvalKey) != 16 {
		log.Fatalf("Expect eval key to be len 16, len=%v, key=%v", len(a.EvalKey), a.EvalKey)
	}

	return a
}

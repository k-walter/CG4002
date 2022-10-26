package common

import (
	"flag"
	"log"
)

type Arg struct {
	EvalAddr    string
	EvalKey     string
	RelayPort   uint
	PynqPort    uint
	ShieldErrNs uint64
	ShootErrNs  uint64
}

func ParseArgs() Arg {
	// Bind a
	a := Arg{}
	flag.StringVar(&a.EvalAddr, "evalAddr", ":8080", "Eval server IP and port")
	flag.StringVar(&a.EvalKey, "evalKey", "PLSPLSPLSPLSWORK", "Symmetric key for eval server")
	flag.UintVar(&a.RelayPort, "relayPort", 8081, "Port to be SSH forwarded to relay")
	flag.UintVar(&a.PynqPort, "pynqPort", 8082, "Port to send/poll signals")
	flag.Uint64Var(&a.ShootErrNs, "ShootErrNs", 200_000_000, "Max difference between when shoot/shot are sent. If too short, will miss shots. If too long, missed shoots will update eval slower")
	flag.Uint64Var(&a.ShieldErrNs, "ShieldErrNs", 100_000_000, "Error between eval's response and current state to reset slow/fast timing")
	flag.Parse()

	// Checks
	if len(a.EvalKey) != 16 {
		log.Fatalf("Expect eval key to be len 16, len=%v, key=%v", len(a.EvalKey), a.EvalKey)
	}

	return a
}

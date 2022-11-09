package common

import (
	"flag"
	"log"
	"time"
)

type Arg struct {
	// Connection
	EvalAddr  string
	EvalKey   string
	RelayPort uint
	PynqPort  uint

	// Visualizer MQTT
	Broker         string
	EventTopic     string
	InFovRespTopic string
	StateTopic     string

	// Health
	HpMax       uint32 // instant rebirth
	ShieldMax   uint32 // per life
	ShieldHpMax uint32
	ShieldTime  time.Duration // reset in next life

	// Damage
	BulletMax  uint32 // unlimited mags, reload only if mag=0
	BulletDmg  uint32
	GrenadeDmg uint32
	GrenadeMax uint32 // per life
	//GrenadeTime time.Duration // to display

	// Tuning
	ShootErr time.Duration
	LPF      float64
}

func ParseArgs() Arg {
	// Bind a
	a := Arg{}
	flag.StringVar(&a.EvalAddr, "evalAddr", ":8080", "Eval server IP and port")
	flag.StringVar(&a.EvalKey, "evalKey", "PLSPLSPLSPLSWORK", "Symmetric key for eval server")
	flag.UintVar(&a.RelayPort, "relayPort", 8081, "Port to be SSH forwarded to relay")
	flag.UintVar(&a.PynqPort, "pynqPort", 8082, "Port to send/poll signals")

	flag.StringVar(&a.Broker, "Broker", "tcp://broker.hivemq.com:1883", "MQTT protocol and URL")
	flag.StringVar(&a.EventTopic, "EventTopic", "cg4002/b7/event", "MQTT topic for events")
	flag.StringVar(&a.InFovRespTopic, "InFovRespTopic", "cg4002/b7/inFovResp", "")
	flag.StringVar(&a.StateTopic, "StateTopic", "cg4002/b7/state", "")

	flag.Float64Var(&a.LPF, "LPF", 0.5, "LPF factor for RTT estimation to eval")

	// Checks
	flag.Parse()
	if len(a.EvalKey) != 16 {
		log.Fatalf("Expect eval key to be len 16, len=%v, key=%v", len(a.EvalKey), a.EvalKey)
	}

	// Health
	a.HpMax = uint32(*flag.Uint64("HpMax", 100, "HP per round"))
	a.ShieldMax = uint32(*flag.Uint64("ShieldMax", 3, "Number of shield per round"))
	a.ShieldHpMax = uint32(*flag.Uint64("ShieldHpMax", 30, "Shield HP when shielding"))
	a.ShieldTime = time.Second * time.Duration(*flag.Uint64("ShieldTime", 10, "Timeout for shield (s)"))

	// Damage
	a.BulletMax = uint32(*flag.Uint64("BulletMax", 6, "Number of bullets per reload/round"))
	a.BulletDmg = uint32(*flag.Uint64("BulletDmg", 10, "Damage inflicted per shot"))
	a.GrenadeDmg = uint32(*flag.Uint64("GrenadeDmg", 30, "Damage inflicted per grenade"))
	a.GrenadeMax = uint32(*flag.Uint64("GrenadeMax", 2, "Number of grenades per round"))

	// Tuning
	a.ShootErr = time.Millisecond * time.Duration(*flag.Uint64("ShootErrMs", 300, "Error for shoot/shot (ms)"))

	return a
}

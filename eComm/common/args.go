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
	MockEval bool
}

func ParseArgs() Arg {
	a := Arg{}

	// Connection
	flag.StringVar(&a.EvalAddr, "evalAddr", "", "Eval server IP and port. Leave empty for mock eval.")
	flag.StringVar(&a.EvalKey, "evalKey", "PLSPLSPLSPLSWORK", "Symmetric key for eval server")
	flag.UintVar(&a.RelayPort, "relayPort", 8081, "Port to be SSH forwarded to relay")
	flag.UintVar(&a.PynqPort, "pynqPort", 8082, "Port to send/poll signals")

	// Visualizer MQTT
	flag.StringVar(&a.Broker, "Broker", "tcp://broker.hivemq.com:1883", "MQTT protocol and URL")
	flag.StringVar(&a.EventTopic, "EventTopic", "cg4002/b7/event", "MQTT topic for events")
	flag.StringVar(&a.InFovRespTopic, "InFovRespTopic", "cg4002/b7/inFovResp", "")
	flag.StringVar(&a.StateTopic, "StateTopic", "cg4002/b7/state", "")

	// Health
	hpMax := flag.Uint64("HpMax", 100, "HP per round")
	shieldMax := flag.Uint64("ShieldMax", 3, "Number of shield per round")
	ShieldHpMax := flag.Uint64("ShieldHpMax", 30, "Shield HP when shielding")
	shieldTime := flag.Uint64("ShieldTime", 10, "Timeout for shield (s)")

	// Damage
	bulletMax := flag.Uint64("BulletMax", 6, "Number of bullets per reload/round")
	bulletDmg := flag.Uint64("BulletDmg", 10, "Damage inflicted per shot")
	grenadeDmg := flag.Uint64("GrenadeDmg", 30, "Damage inflicted per grenade")
	grenadeMax := flag.Uint64("GrenadeMax", 2, "Number of grenades per round")

	// Tuning
	shootErr := flag.Uint64("ShootErrMs", 300, "Error for shoot/shot (ms)")
	flag.Float64Var(&a.LPF, "LPF", 0.5, "LPF factor for RTT estimation to eval")

	// Bind and check
	flag.Parse()
	if len(a.EvalKey) != 16 {
		log.Fatalf("Expect eval key to be len 16, len=%v, key=%v", len(a.EvalKey), a.EvalKey)
	}

	// Set vars
	a.HpMax = uint32(*hpMax)
	a.ShieldMax = uint32(*shieldMax)
	a.ShieldHpMax = uint32(*ShieldHpMax)
	a.ShieldTime = time.Second * time.Duration(*shieldTime)
	a.BulletMax = uint32(*bulletMax)
	a.BulletDmg = uint32(*bulletDmg)
	a.GrenadeDmg = uint32(*grenadeDmg)
	a.GrenadeMax = uint32(*grenadeMax)
	a.ShootErr = time.Millisecond * time.Duration(*shootErr)
	a.MockEval = a.EvalAddr == ""

	return a
}

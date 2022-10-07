package common

import (
	pb "cg4002/protos"
	"time"
)

const (
	NPlayer     = 2
	ShootErrNs  = 100_000_000 // 100ms
	ShieldErrNs = 100_000_000 // 100ms

	// Health
	HpMax       = 100 // instant rebirth
	ShieldMax   = 3   // per life
	ShieldHpMax = 30
	ShieldNs    = uint64(10_000_000_000) // 10s, carry over to next life

	// Damage
	BulletMax   = 6 // unlimited mags, reload only if mag=0
	BulletDmg   = 10
	GrenadeDmg  = 10
	GrenadeMax  = 2               // per life
	GrenadeSecs = 2 * time.Second // to damage and display
)

type PlayerImpl struct {
	Shoot []uint64
	Shot  []uint64
}

type PlayerState struct {
	*pb.PlayerState // DO NOT pass ptr outside engine, is race condition
	PlayerImpl
}

func NewState() PlayerState {
	return PlayerState{
		PlayerState: &pb.PlayerState{
			Hp:           HpMax,
			Action:       pb.Action_none,
			Bullets:      BulletMax,
			Grenades:     GrenadeMax,
			ShieldTime:   0,
			ShieldHealth: 0,
			NumDeaths:    0,
			NumShield:    ShieldMax,
		},
		PlayerImpl: PlayerImpl{},
	}
}

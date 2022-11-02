package common

import (
	pb "cg4002/protos"
	"time"
)

const (
	NPlayer = 2

	// Health
	HpMax       = 100 // instant rebirth
	ShieldMax   = 3   // per life
	ShieldHpMax = 30
	ShieldNs    = 10 * time.Second // reset in next life

	// Damage
	BulletMax   = 6 // unlimited mags, reload only if mag=0
	BulletDmg   = 10
	GrenadeDmg  = 30
	GrenadeMax  = 2               // per life
	GrenadeSecs = 2 * time.Second // to display

	ShieldRst = 0
)

type PlayerImpl struct {
	ShieldExpireNs uint64
	Shoot          map[uint32]struct{}
	Shot           map[uint32]struct{}
}

type PlayerState struct {
	*pb.PlayerState // DO NOT pass ptr outside engine, is race condition
	PlayerImpl
}

func NewState() *pb.State {
	return &pb.State{
		P1: NewPlayerState(),
		P2: NewPlayerState(),
	}
}

func NewPlayerState() *pb.PlayerState {
	return &pb.PlayerState{
		Hp:           HpMax,
		Action:       pb.Action_none,
		Bullets:      BulletMax,
		Grenades:     GrenadeMax,
		ShieldTime:   0,
		ShieldHealth: 0,
		NumDeaths:    0,
		NumShield:    ShieldMax,
	}
}

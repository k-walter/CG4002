package common

import (
	pb "cg4002/protos"
	"time"
)

const (
	NPlayer = 2
)

type PlayerImpl struct {
	ShieldExpireNs uint64
	Shoot          map[uint32]time.Time
	Shot           map[uint32]time.Time
}

type PlayerState struct {
	*pb.PlayerState // DO NOT pass ptr outside engine, is race condition
	PlayerImpl
}

func NewState(a *Arg) *pb.State {
	return &pb.State{
		P1: NewPlayerState(a),
		P2: NewPlayerState(a),
	}
}

func NewPlayerState(a *Arg) *pb.PlayerState {
	return &pb.PlayerState{
		Hp:           a.HpMax,
		Action:       pb.Action_none,
		Bullets:      a.BulletMax,
		Grenades:     a.GrenadeMax,
		ShieldTime:   0,
		ShieldHealth: 0,
		NumDeaths:    0,
		NumShield:    a.ShieldMax,
	}
}

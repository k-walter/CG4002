package engine

import (
	cmn "cg4002/eComm/common"
	"log"
)

// Match then clear preceding shoot/shot, else no updates
// Goal: update eval when match (should clear preceding shoots), or when timeout (shootTimeout clears preceding shoots)
func matchShot(shootID uint32, s *PlayerImpl) bool {
	// Match shoot and shot
	if _, found := s.shoot[shootID]; !found {
		return false
	}
	if _, found := s.shot[shootID]; !found {
		return false
	}

	// inflict dmg
	log.Println("matched shootID =", shootID)

	return true
}

func inflict(p *PlayerImpl, dmg uint32, a *cmn.Arg) {
	if p.Hp+p.ShieldHealth <= dmg {
		// Die & revive
		p.NumDeaths += 1
		p.Hp = a.HpMax
		p.ShieldHealth = 0
		p.Grenades = a.GrenadeMax
		p.NumShield = a.ShieldMax
		p.Bullets = a.BulletMax

		// RULE reset shield cooldown
		if p.shieldExpiry.After(cmn.GameTime) {
			handleShieldAvail(p)
		}

	} else if p.ShieldHealth <= dmg {
		// Dmg shield+health
		p.Hp -= dmg - p.ShieldHealth
		p.ShieldHealth = 0
	} else {
		// Dmg shield
		p.ShieldHealth -= dmg
	}
}

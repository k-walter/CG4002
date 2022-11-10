package engine

import "cg4002/eComm/common"

func handleShootTimeout(p *PlayerImpl) {
	if p.fsm != shotBullet {
		return
	}
	p.fsm = done
	if !p.shootTimeout.Stop() {
		common.Drain(p.shootTimeout.C)
	}
}

func handleShieldAvail(p *PlayerImpl) {
	p.ShieldHealth = 0
	p.ShieldTime = 0
	p.shieldExpiry = common.GameTime
	if !p.shieldTimeout.Stop() {
		common.Drain(p.shieldTimeout.C)
	}
}

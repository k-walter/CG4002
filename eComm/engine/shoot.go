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
		log.Printf("could not find shooter, %v\n", s.shoot)
		return false
	}
	if _, found := s.shot[shootID]; !found {
		log.Printf("could not find victim, %v\n", s.shot)
		return false
	}

	// inflict dmg
	log.Println("matched shootID =", shootID)
	inflict(s, cmn.BulletDmg)

	return true
}

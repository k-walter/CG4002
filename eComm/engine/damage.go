package engine

import (
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

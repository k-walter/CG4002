package common

import (
	"log"
	"sync"
	"time"
)

type Number interface {
	uint64
}

func Max[T Number](i, j T) T {
	if i > j {
		return i
	}
	return j
}

func Min[T Number](i, j T) T {
	if i < j {
		return i
	}
	return j
}

func AbsDiff[T Number](i, j T) T {
	return Max(i, j) - Min(i, j)
}

func EXIT_UNLESS(b bool) {
	if !b {
		log.Panicln()
	}
}

// Returns index if found, else -1
func BinarySearch[T Number](a []T, v T) int {
	for l, r := 0, len(a); l < r; {
		m := (l + r) >> 1
		if a[m] == v {
			return m
		}
		if a[m] < v {
			l = m + 1
		} else {
			r = m
		}
	}
	return -1
}

func Do(wg *sync.WaitGroup, f func()) {
	f()
	wg.Done()
}

func Drain(c <-chan time.Time) {
	for ok := true; ok; {
		select {
		case _, ok = <-c:
		default:
			return
		}
	}
}

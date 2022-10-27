package common

import "time"

var GameTime time.Time
var RndNum uint32

func MakeClock() {
	RndNum = 1
	GameTime = time.Now()
}

// TimeToNs converts monotonic clk --> duration --(NsToTime)--> monotonic clk
func TimeToNs(t time.Time) uint64 {
	return uint64(t.Sub(GameTime).Nanoseconds())
}

func NsToTime(ns uint64) time.Time {
	return GameTime.Add(time.Duration(ns) * time.Nanosecond)
}

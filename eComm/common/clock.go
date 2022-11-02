package common

import "time"

type RoundT uint32

var GameTime time.Time = time.Now()

// TimeToNs converts monotonic clk --> duration --(NsToTime)--> monotonic clk
func TimeToNs(t time.Time) uint64 {
	return uint64(t.Sub(GameTime).Nanoseconds())
}

func NsToTime(ns uint64) time.Time {
	return GameTime.Add(time.Duration(ns) * time.Nanosecond)
}

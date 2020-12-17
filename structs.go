package profiling

import "time"

type Profiling struct {
	printMode PrintMode
}

type PrintMode int

type gcSummary struct {
	NumGC      int64
	PauseTotal time.Duration
	LastPause  time.Duration
	PauseAvg   time.Duration
	Overhead   float64
	Alloc      uint64
	AllocRate  uint64
	Sys        uint64
	Histogram1 time.Duration
	Histogram2 time.Duration
	Histogram3 time.Duration
}

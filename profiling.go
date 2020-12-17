package profiling

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"time"
)

func newProfiling(pm PrintMode) *Profiling {
	return &Profiling{printMode: pm}
}

func SetPrintMode(d PrintMode) {
	profiling.printMode = d
}

// GoRoutine prints go routines information
func GoRoutine(w io.Writer) error {
	return pprof.Lookup(pprofGoRoutine).WriteTo(w, int(profiling.printMode))
}

// ThreadCreate prints go thread create information
func ThreadCreate(w io.Writer) error {
	return pprof.Lookup(pprofThreadCreate).WriteTo(w, int(profiling.printMode))
}

// Heap prints heap information
func Heap(w io.Writer) error {
	return pprof.Lookup(pprofHeap).WriteTo(w, int(profiling.printMode))
}

// Allocs prints allocations information
func Allocs(w io.Writer) error {
	return pprof.Lookup(pprofAllocs).WriteTo(w, int(profiling.printMode))
}

// Block prints block information
func Block(rate int, w io.Writer) error {
	defer runtime.SetBlockProfileRate(0)
	runtime.SetBlockProfileRate(rate)
	return pprof.Lookup(pprofBlock).WriteTo(w, int(profiling.printMode))
}

// Mutex prints mutex information
func Mutex(rate int, w io.Writer) error {
	defer runtime.SetMutexProfileFraction(0)
	runtime.SetMutexProfileFraction(rate)
	return pprof.Lookup(pprofMutex).WriteTo(w, int(profiling.printMode))
}

// CPU check it using: go tool pprof <file>
func CPU(d time.Duration, w io.Writer) (err error) {
	if err = pprof.StartCPUProfile(w); err != nil {
		return err
	}

	time.Sleep(d)

	pprof.StopCPUProfile()

	return nil
}

// Memory prints memory information
func Memory(w io.Writer) (err error) {
	runtime.GC()
	return pprof.WriteHeapProfile(w)
}

// GC prints garbage collection information
func GC(w io.Writer) error {
	startTime := time.Now()

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	gcStats := &debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(gcStats)

	printGC(startTime, memStats, gcStats, w)

	return nil
}

// Trace prints a program trace
func Trace(d time.Duration, w io.Writer) error {
	if err := trace.Start(w); err != nil {
		return err
	}

	time.Sleep(d)

	runtime.StopTrace()

	return nil
}

// Symbol prints functions information
func Symbol(words []string, w io.Writer) error {
	for _, word := range words {
		pc, _ := strconv.ParseUint(word, 0, 64)
		if pc != 0 {
			f := runtime.FuncForPC(uintptr(pc))
			if f != nil {
				fmt.Fprintf(w, "%#x %s\n", pc, f.Name())
			}
		}
	}

	return nil
}

func printGC(startTime time.Time, memStats *runtime.MemStats, gcstats *debug.GCStats, w io.Writer) {
	switch gcstats.NumGC > 0 {
	case true:
		elapsed := time.Now().Sub(startTime)

		summary := gcSummary{
			NumGC:      gcstats.NumGC,
			PauseTotal: gcstats.PauseTotal,
			LastPause:  gcstats.Pause[0],
			PauseAvg:   avg(gcstats.Pause),
			Overhead:   float64(gcstats.PauseTotal) / float64(elapsed) * 100,
			Alloc:      memStats.Alloc,
			Sys:        memStats.Sys,
			AllocRate:  uint64(float64(memStats.TotalAlloc) / elapsed.Seconds()),
			Histogram1: gcstats.PauseQuantiles[94],
			Histogram2: gcstats.PauseQuantiles[98],
			Histogram3: gcstats.PauseQuantiles[99],
		}

		fmt.Fprintf(w, "NumGC:%d Pause:%s Pause(Avg):%s Overhead:%3.2f%% Alloc:%s Sys:%s Alloc(Rate):%s/s Histogram:%s %s %s \n",
			summary.NumGC,
			toS(summary.LastPause),
			toS(summary.PauseAvg),
			summary.Overhead,
			toH(summary.Alloc),
			toH(summary.AllocRate),
			toH(summary.Sys),
			toS(summary.Histogram1),
			toS(summary.Histogram2),
			toS(summary.Histogram3))
	case false:
		// while GC has disabled
		elapsed := time.Now().Sub(startTime)

		summary := gcSummary{
			Alloc:     memStats.Alloc,
			Sys:       memStats.Sys,
			AllocRate: uint64(float64(memStats.TotalAlloc) / elapsed.Seconds()),
		}

		fmt.Fprintf(w, "Alloc:%s Sys:%s Alloc(Rate):%s/s\n",
			toH(summary.Alloc),
			toH(summary.Sys),
			toH(summary.AllocRate))
	}
}

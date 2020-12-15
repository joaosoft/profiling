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
func Block(w io.Writer) error {
	return pprof.Lookup(pprofBlock).WriteTo(w, int(profiling.printMode))
}

// Mutex prints mutex information
func Mutex(w io.Writer) error {
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
		lastPause := gcstats.Pause[0]
		elapsed := time.Now().Sub(startTime)
		overhead := float64(gcstats.PauseTotal) / float64(elapsed) * 100
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "NumGC:%d Pause:%s Pause(Avg):%s Overhead:%3.2f%% Alloc:%s Sys:%s Alloc(Rate):%s/s Histogram:%s %s %s \n",
			gcstats.NumGC,
			toS(lastPause),
			toS(avg(gcstats.Pause)),
			overhead,
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)),
			toS(gcstats.PauseQuantiles[94]),
			toS(gcstats.PauseQuantiles[98]),
			toS(gcstats.PauseQuantiles[99]))
	case false:
		// while GC has disabled
		elapsed := time.Now().Sub(startTime)
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "Alloc:%s Sys:%s Alloc(Rate):%s/s\n",
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)))
	}
}

func avg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

// format bytes number friendly
func toH(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}

// short string format
func toS(d time.Duration) string {

	u := uint64(d)
	if u < uint64(time.Second) {
		switch {
		case u == 0:
			return "0"
		case u < uint64(time.Microsecond):
			return fmt.Sprintf("%.2fns", float64(u))
		case u < uint64(time.Millisecond):
			return fmt.Sprintf("%.2fus", float64(u)/1000)
		default:
			return fmt.Sprintf("%.2fms", float64(u)/1000/1000)
		}
	} else {
		switch {
		case u < uint64(time.Minute):
			return fmt.Sprintf("%.2fs", float64(u)/1000/1000/1000)
		case u < uint64(time.Hour):
			return fmt.Sprintf("%.2fm", float64(u)/1000/1000/1000/60)
		default:
			return fmt.Sprintf("%.2fh", float64(u)/1000/1000/1000/60/60)
		}
	}
}

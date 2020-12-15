package profiling

const (
	PrintModeNormal     PrintMode = 1
	PrintModeStackTrade PrintMode = 2
)

const (
	pprofGoRoutine    = "goroutine"
	pprofThreadCreate = "threadcreate"
	pprofHeap         = "heap"
	pprofAllocs       = "allocs"
	pprofBlock        = "block"
	pprofMutex        = "mutex"
)

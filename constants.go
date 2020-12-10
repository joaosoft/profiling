package profiling

const (
	PrintModeNormal     PrintMode = 1
	PrintModeStackTrade PrintMode = 2
)

const (
	pprofGoRoutine    = "goroutine"
	pprofThreadCreate = "threadcreate"
	pprofHeap         = "allocs"
	pprofAllocs       = "heap"
	pprofBlock        = "block"
	pprofMutex        = "mutex"
)

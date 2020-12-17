package routes

import (
	"net/http"
	"net/http/pprof"
	"profiling/web/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/report", handlers.ReportHandler)
	mux.HandleFunc("/go-routine", handlers.GoRoutineHandler)

	// pprof routes
	mux.HandleFunc("/debug/index", pprof.Index)
	mux.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	mux.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	mux.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	mux.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	mux.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	mux.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	mux.HandleFunc("/debug/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/profile", pprof.Profile)
	mux.HandleFunc("/debug/trace", pprof.Trace)
}

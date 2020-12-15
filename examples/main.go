package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"profiling"
	"time"
)

const (
	generatedFolder = "./generated"
)

var (
	quit = make(chan bool)
)

func init() {
	profiling.SetPrintMode(profiling.PrintModeNormal)
	initProcesses(quit)
}

func initProcesses(quit chan bool) {
	for i := 0; i < 10; i++ {
		for true {
			go func() {
				select {
				case <-quit:
					fmt.Println("received shutdown signal")
				}
			}()
		}
	} // profiling over http using pprof
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/profile", pprof.Profile)
	mux.HandleFunc("/dummy", dummyHandler)
	go http.ListenAndServe(":7777", mux)

}

func main() {
	var err error
	w := &bytes.Buffer{}

	// GoRoutine
	fmt.Fprint(w, "\n\n:: Executing: GoRoutine\n")
	if err = profiling.GoRoutine(w); err != nil {
		panic(err)
	}

	// ThreadCreate
	fmt.Fprint(w, "\n\n:: Executing: ThreadCreate\n")
	if err = profiling.ThreadCreate(w); err != nil {
		panic(err)
	}

	// Heap
	fmt.Fprint(w, "\n\n:: Executing: Heap\n")
	if err = profiling.Heap(w); err != nil {
		panic(err)
	}

	// Allocs
	fmt.Fprint(w, "\n\n:: Executing: Allocs\n")
	if err = profiling.Allocs(w); err != nil {
		panic(err)
	}

	// Block
	fmt.Fprint(w, "\n\n:: Executing: Block\n")
	if err = profiling.Block(w); err != nil {
		panic(err)
	}

	// Mutex
	fmt.Fprint(w, "\n\n:: Executing: Mutex\n")
	if err = profiling.Mutex(w); err != nil {
		panic(err)
	}

	// CPU
	fmt.Println(":: Executing: CPU")
	fileName := fmt.Sprintf("%s/cpu-%d.pprof", generatedFolder, os.Getpid())
	var file *os.File
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	if err = profiling.CPU(20*time.Second, file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool pprof %s\n", fileName)

	// Memory
	fmt.Println(":: Executing: Memory")
	fileName = fmt.Sprintf("%s/mem-%d.memprof", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	if err = profiling.Memory(file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool pprof %s\n", fileName)

	fmt.Fprint(w, "\n\n:: Executing: GC\n")
	if err = profiling.GC(w); err != nil {
		panic(err)
	}

	// Print w to stdout
	fmt.Println(w.String())

	fmt.Println("quit goroutines")
	quit <- true
}

func dummyHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("> executing dummy ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

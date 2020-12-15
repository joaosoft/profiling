package main

import (
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
	//initProcesses(quit)
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
	}

	// profiling over http using pprof
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/profile", pprof.Profile)
	mux.HandleFunc("/dummy", dummyHandler)
	go http.ListenAndServe(":7777", mux)
}

func main() {
	var err error
	var file *os.File
	var fileName string

	// GoRoutine
	fmt.Println(":: Executing: GoRoutine")
	fileName = fmt.Sprintf("%s/go-routine-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.GoRoutine(file); err != nil {
		panic(err)
	}

	// ThreadCreate
	fmt.Println(":: Executing: ThreadCreate")
	fileName = fmt.Sprintf("%s/thread-create-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.ThreadCreate(file); err != nil {
		panic(err)
	}

	// Heap
	fmt.Println(":: Executing: Heap")
	fileName = fmt.Sprintf("%s/heap-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.Heap(file); err != nil {
		panic(err)
	}

	// Allocs
	fmt.Println(":: Executing: Allocs")
	fileName = fmt.Sprintf("%s/allocs-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.Allocs(file); err != nil {
		panic(err)
	}

	// Block
	fmt.Println(":: Executing: Block")
	fileName = fmt.Sprintf("%s/block-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.Block(file); err != nil {
		panic(err)
	}

	// Mutex
	fmt.Println(":: Executing: Mutex")
	fileName = fmt.Sprintf("%s/mutex-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.Mutex(file); err != nil {
		panic(err)
	}

	// Trace
	fmt.Println(":: Executing: Trace")
	fileName = fmt.Sprintf("%s/trace-%d.pprof", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.Trace(20*time.Second, file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool trace %s\n", fileName)

	// Symbol
	fmt.Println(":: Executing: Symbol")
	fileName = fmt.Sprintf("%s/symbol-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.Symbol([]string{"main"}, file); err != nil {
		panic(err)
	}

	// CPU
	fmt.Println(":: Executing: CPU")
	fileName = fmt.Sprintf("%s/cpu-%d.pprof", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
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
	defer file.Close()
	if err = profiling.Memory(file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool pprof %s\n", fileName)

	// GC
	fmt.Println(":: Executing: GC")
	fileName = fmt.Sprintf("%s/gc-%d.text", generatedFolder, os.Getpid())
	file, err = os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = profiling.GC(file); err != nil {
		panic(err)
	}

	fmt.Println("quit goroutines")
	quit <- true
}

func dummyHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("> executing dummy")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

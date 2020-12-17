package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"profiling"
	"syscall"
	"time"
)

const (
	generatedFolder   = "./generated"
	httpWebServerPort = 7777
	numGoRoutines     = 50
)

func init() {
	profiling.SetPrintMode(profiling.PrintModeNormal)
}

func main() {
	startWebServer(httpWebServerPort)
	stopChan := make(chan bool)
	startDummyProcesses(stopChan)
	startProfileTools()

	close(stopChan) // tell it to stop
	<-stopChan      // wait for it to have stopped

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	fmt.Println("waiting for term command")
	<-termChan
}

func startWebServer(port int) {
	mux := http.NewServeMux()

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

	// dummy route
	mux.HandleFunc("/dummy", dummyHandler)

	fmt.Printf("web server started at http://localhost:%d/debug/index\n", port)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func dummyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var data []int
	for i := 0; i < 100000; i++ {
		if i%1000 == 0 {
			data = append(data, i)
		}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bytes)
}

func startDummyProcesses(stopChan chan bool) {
	for i := 0; i <= numGoRoutines; i++ {
		go func(id int) {
			for true {
				select {
				default:
					_, err := http.Get(fmt.Sprintf("http://localhost:%d/dummy", httpWebServerPort))
					if err != nil {
						panic(err)
					}
					time.Sleep(time.Second)
				case <-stopChan:
					fmt.Printf("\nstopping go routine %d", id)
					// stop
					return
				}
			}
		}(i)
	}
}

func startProfileTools() {
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
	if err = profiling.Block(100, file); err != nil {
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
	if err = profiling.Mutex(100, file); err != nil {
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
	if err = profiling.Trace(30*time.Second, file); err != nil {
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
	if err = profiling.CPU(30*time.Second, file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool pprof %s\n", fileName)

	// show on prof UI
	fmt.Println("showing on pprof UI")
	showPprofUI(fileName)

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

	// show on prof UI
	fmt.Println("showing on pprof UI")
	showPprofUI(fileName)

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

	fmt.Println("done")
}

func showPprofUI(fileName string) {
	cmd := exec.Command("go", "tool", "pprof", "-http=:", fileName)

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("cannot start pprof UI: %v", err))
	}
}

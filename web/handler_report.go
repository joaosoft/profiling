package web

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"profiling"
	"time"
)

func reportHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	stopChan := make(chan bool)
	startDummyProcesses(httpWebServerPort, numGoRoutines, stopChan)
	startProfileTools(outputFolder)

	close(stopChan) // tell it to stop
	<-stopChan      // wait for it to have stopped

	w.Write([]byte("Report generated!"))
}

func startDummyProcesses(port, numGoRoutines int, stopChan chan bool) {
	for i := 0; i <= numGoRoutines; i++ {
		go func(id int) {
			for true {
				select {
				default:
					url := fmt.Sprintf("http://localhost:%d/go-routine", port)
					_, err := http.Get(url)
					if err != nil {
						panic(err)
					}
					time.Sleep(time.Second)
				case <-stopChan:
					// stop
					return
				}
			}
		}(i)
	}
}

func createFile(folder, name, extension string) *os.File {
	fileName := fmt.Sprintf("%s/%s-%d.%s", folder, name, os.Getpid(), extension)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	return file
}

func startProfileTools(outputFolder string) {
	var err error
	var file *os.File

	// GoRoutine
	fmt.Println(":: Executing: GoRoutine")
	file = createFile(outputFolder, "go-routine", "text")
	defer file.Close()
	if err = profiling.GoRoutine(file); err != nil {
		panic(err)
	}

	// ThreadCreate
	fmt.Println(":: Executing: ThreadCreate")
	file = createFile(outputFolder, "thread-create", "text")
	defer file.Close()
	if err = profiling.ThreadCreate(file); err != nil {
		panic(err)
	}

	// Heap
	fmt.Println(":: Executing: Heap")
	file = createFile(outputFolder, "heap", "text")
	defer file.Close()
	if err = profiling.Heap(file); err != nil {
		panic(err)
	}

	// Allocs
	fmt.Println(":: Executing: Allocs")
	file = createFile(outputFolder, "allocs", "text")
	defer file.Close()
	if err = profiling.Allocs(file); err != nil {
		panic(err)
	}

	// Block
	fmt.Println(":: Executing: Block")
	file = createFile(outputFolder, "block", "text")
	defer file.Close()
	if err = profiling.Block(100, file); err != nil {
		panic(err)
	}

	// Mutex
	fmt.Println(":: Executing: Mutex")
	file = createFile(outputFolder, "mutex", "text")
	defer file.Close()
	if err = profiling.Mutex(100, file); err != nil {
		panic(err)
	}

	// Trace
	fmt.Println(":: Executing: Trace")
	file = createFile(outputFolder, "trace", "pprof")
	defer file.Close()
	if err = profiling.Trace(30*time.Second, file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool trace %s\n", file.Name())

	// show on prof UI
	fmt.Println("showing on UI")
	showGoToolUI(goToolCmdTrace, file.Name())

	// Symbol
	fmt.Println(":: Executing: Symbol")
	file = createFile(outputFolder, "symbol", "text")
	defer file.Close()
	if err = profiling.Symbol([]string{"main"}, file); err != nil {
		panic(err)
	}

	// CPU
	fmt.Println(":: Executing: CPU")
	file = createFile(outputFolder, "cpu", "pprof")
	defer file.Close()
	if err = profiling.CPU(30*time.Second, file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool pprof %s\n", file.Name())

	// show on prof UI
	fmt.Println("showing on UI")
	showGoToolUI(goToolCmdProf, file.Name())

	// Memory
	fmt.Println(":: Executing: Memory")
	file = createFile(outputFolder, "mem", "memprof")
	defer file.Close()
	if err = profiling.Memory(file); err != nil {
		panic(err)
	}
	fmt.Printf("Now you can use the command: go tool pprof %s\n", file.Name())

	// show on prof UI
	fmt.Println("showing on UI")
	showGoToolUI(goToolCmdProf, file.Name())

	// GC
	fmt.Println(":: Executing: GC")
	file = createFile(outputFolder, "gc", "text")
	defer file.Close()
	if err = profiling.GC(file); err != nil {
		panic(err)
	}

	fmt.Println("done")
}

func showGoToolUI(command, fileName string) {
	cmd := exec.Command("go", "tool", command, "-http=:", fileName)

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("cannot start pprof UI: %v", err))
	}
}

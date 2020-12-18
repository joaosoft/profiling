package web

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

func init() {
	mux = http.NewServeMux()
}

func Start() {
	registerRoutes(mux)

	url := fmt.Sprintf("http://localhost:%d", httpWebServerPort)
	fmt.Printf("web server started at %s\n", url)

	if err := showUrl(url); err != nil {
		log.Println(fmt.Sprintf("cannot start browser: %v", err))
	}

	http.ListenAndServe(fmt.Sprintf(":%d", httpWebServerPort), mux)
}

func showUrl(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

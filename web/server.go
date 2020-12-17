package web

import (
	"fmt"
	"net/http"
	"profiling/web/routes"
)

func init() {
	mux = http.NewServeMux()
}

func Start() {
	routes.RegisterRoutes(mux)

	fmt.Printf("web server started at http://localhost:%d\n", HttpWebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%d", HttpWebServerPort), mux)
}


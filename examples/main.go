package main

import (
	"github.com/joaosoft/profiling"
	"github.com/joaosoft/profiling/web"
)

func init() {
	profiling.SetPrintMode(profiling.PrintModeNormal)
}

func main() {
	web.Start()
}

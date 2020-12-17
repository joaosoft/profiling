package main

import (
	"profiling"
	"profiling/web"
)

func init() {
	profiling.SetPrintMode(profiling.PrintModeNormal)
}

func main() {
	web.Start()
}

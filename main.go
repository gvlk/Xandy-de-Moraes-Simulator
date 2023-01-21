package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	win, _ := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Xandy D'Moraes Simulator",
			Bounds: pixel.R(0, 0, 1024, 768),
		})
	gameControl(win)
}

func main() {
	pixelgl.Run(run)
}

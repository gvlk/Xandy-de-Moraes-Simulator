package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Xandy D'Moraes Simulator",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	gamecontroller.start(gamecontroller{cfg: cfg})
}

func main() {
	pixelgl.Run(run)
}

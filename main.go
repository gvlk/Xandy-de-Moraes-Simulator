package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	sWidth  = 1920.0
	sHeight = sWidth * (9.0 / 16.0)
	icon    = loadPicture("images/icon.png")
)

func run() {
	win, _ := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:       "Xandy D'Moraes Simulator",
			Bounds:      pixel.R(0, 0, sWidth, sHeight),
			Undecorated: true,
			Icon:        []pixel.Picture{icon},
		})
	gameControl(win)
}

func main() {
	pixelgl.Run(run)
}

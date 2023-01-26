package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	sWidth  = 1920.0
	sHeight = sWidth * (9.0 / 16.0)
)

var sCenter = pixel.V(sWidth/2, sHeight/2)

func run() {
	win, _ := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:       "Xandy de Moraes Simulator",
			Bounds:      pixel.R(0, 0, sWidth, sHeight),
			Undecorated: true,
			Icon:        []pixel.Picture{loadPicture("images/icon.png")},
		})
	gameControl(win)
}

func main() {
	pixelgl.Run(run)
}

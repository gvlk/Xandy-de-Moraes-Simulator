package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var debugAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

var surface = &pixelgl.Window{}

func setDebugWindow(win *pixelgl.Window) {
	surface = win
}

func displayDebug(info string, y int) {

	txtBox := text.New(pixel.V(10, surface.Bounds().Max.Y-float64(y*15)), debugAtlas)
	txtBox.Color = colornames.Black
	txtBox.WriteString(info)
	txtBox.Draw(surface, pixel.IM)
}

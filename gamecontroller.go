package main

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	//"github.com/faiface/pixel"
	//"github.com/faiface/pixel/pixelgl"
)

type gamecontroller struct {
	cfg pixelgl.WindowConfig
}

func (gc gamecontroller) start() {
	cfg := gc.cfg
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	for !win.Closed() {
		win.Update()
	}

}

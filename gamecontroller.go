package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"time"
)

func start(wc pixelgl.WindowConfig) {
	win, err := pixelgl.NewWindow(wc)
	if err != nil {
		panic(err)
	}

	face, err := loadTTF("fonts/intuitive.ttf", 80)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(100, 500), atlas)
	txt.Color = colornames.Black

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {

		txt.Clear()
		fmt.Fprintf(txt, "%d:%d:%d", time.Now().Hour(), time.Now().Minute(), time.Now().Second())

		win.Clear(colornames.Whitesmoke)
		txt.Draw(win, pixel.IM)
		win.Update()
		<-fps
	}
}

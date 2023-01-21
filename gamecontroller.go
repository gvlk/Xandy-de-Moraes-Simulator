package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
	"unicode"
)

func gameControl(win *pixelgl.Window) {

	face, err := loadTTF("fonts/intuitive.ttf", 12)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII, text.RangeTable(unicode.Latin))
	txt := text.New(pixel.V(50, 700), atlas)
	txt.Color = colornames.Black
	caso := "O réu Eurico Nunes é acusado pelo homicídio de Flávia Medeiros.\nA promotoria argumenta que eles eram casados, mas estavam em processo de separação judicial e já não viviam mais juntos, enquanto Eurico não aceitava muito bem o término.\nA defesa argumenta que Eurico mudou a abordagem quando conheceu uma nova mulher, que testemunhou no tribunal e afirmou categoricamente que viu Eurico na casa."
	fmt.Fprintln(txt, caso)

	// Create an imdraw
	imd := imdraw.New(nil)

	// Set the color of the rectangle
	imd.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// Draw the rectangle around the text
	imd.Push(txt.Bounds().Min, txt.Bounds().Max)
	imd.Rectangle(5)

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {

		win.Clear(colornames.Whitesmoke)
		txt.Draw(win, pixel.IM)
		imd.Draw(win)
		win.Update()
		<-fps
	}
}

package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

func gameControl(win *pixelgl.Window) {

	setDebugWindow(win)

	backgroundPic := loadPicture("images/sprites/cenario01.png")
	backgroundSpr := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())

	casoTexto, casoRect := makeTextBox("txt/cases.txt", pixel.V(500, 300), pixel.V(1000, 530))
	clickTexts1, clickRect1 := makeMultiTextBox("txt/teste.txt", pixel.V(70, 100), pixel.V(160, 250))
	clickTexts2, clickRect2 := makeMultiTextBox("txt/teste.txt", pixel.V(1200, 100), pixel.V(1290, 250))

	fps := time.Tick(time.Second / 60)
	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	// main loop
	for !win.Closed() {

		mousePos := win.MousePosition()

		//events

		//tick
		win.Clear(colornames.Whitesmoke)
		backgroundSpr.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		displayDebug("MousePos: "+mousePos.String(), 0)

		casoTexto.Draw(win, pixel.IM)
		casoRect.Draw(win)
		for _, clickText := range clickTexts1 {
			clickText.Draw(win, pixel.IM)
		}
		for _, clickText := range clickTexts2 {
			clickText.Draw(win, pixel.IM)
		}
		clickRect1.Draw(win)
		clickRect2.Draw(win)

		//hover
		for _, clickText := range clickTexts1 {
			if clickText.Bounds().Contains(mousePos) {
				imd.Clear()
				imd.Push(clickText.Bounds().Min, clickText.Bounds().Max)
				imd.Rectangle(0)
				imd.Draw(win)
			}
		}
		for _, clickText := range clickTexts2 {
			if clickText.Bounds().Contains(mousePos) {
				imd.Clear()
				imd.Push(clickText.Bounds().Min, clickText.Bounds().Max)
				imd.Rectangle(0)
				imd.Draw(win)
			}
		}

		win.Update()
		<-fps
	}
}

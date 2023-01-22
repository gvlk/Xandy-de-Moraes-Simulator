package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

func gameControl(win *pixelgl.Window) {

	setDebugWindow(win)

	//teste := pixelgl.NewCanvas(win.Bounds())
	//teste.SetFragmentShader(fragmentShader)

	backgroundPic := loadPicture("images/sprites/cenario01.png")
	backgroundSpr := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())

	casoTexto, casoRect := makeTextBox("txt/cases.txt", pixel.V(500, 300), pixel.V(1000, 530))
	clickTexts1, clickRect1 := makeMultiClickTextBox("txt/teste.txt", pixel.V(70, 100), pixel.V(160, 250))
	clickTexts2, clickRect2 := makeMultiClickTextBox("txt/teste.txt", pixel.V(1200, 100), pixel.V(1290, 250))

	fps := time.Tick(time.Second / 60)
	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	// main loop
	for !win.Closed() {

		mousePos := win.MousePosition()

		//events
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			for i := range clickTexts1 {
				if clickTexts1[i].txtBox.Bounds().Contains(mousePos) {
					clickTexts1[i].click = !clickTexts1[i].click
				}
			}
			for i := range clickTexts2 {
				if clickTexts2[i].txtBox.Bounds().Contains(mousePos) {
					clickTexts2[i].click = !clickTexts2[i].click
				}
			}
		}

		//tick
		win.Clear(colornames.Whitesmoke)
		backgroundSpr.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		displayDebug("MousePos = "+mousePos.String(), 0)
		displayDebug("clickTexts1[0].click = "+strconv.FormatBool(clickTexts1[0].click), 1)

		casoTexto.Draw(win, pixel.IM)
		casoRect.Draw(win)
		for _, clickText := range clickTexts1 {
			clickText.txtBox.Draw(win, pixel.IM)
		}
		for _, clickText := range clickTexts2 {
			clickText.txtBox.Draw(win, pixel.IM)
		}
		clickRect1.Draw(win)
		clickRect2.Draw(win)

		//hover
		for _, clickText := range clickTexts1 {
			if clickText.txtBox.Bounds().Contains(mousePos) || clickText.click {
				imd.Clear()
				imd.Push(clickText.txtBox.Bounds().Min, clickText.txtBox.Bounds().Max)
				imd.Rectangle(0)
				imd.Draw(win)
			}
		}
		for _, clickText := range clickTexts2 {
			if clickText.txtBox.Bounds().Contains(mousePos) || clickText.click {
				imd.Clear()
				imd.Push(clickText.txtBox.Bounds().Min, clickText.txtBox.Bounds().Max)
				imd.Rectangle(0)
				imd.Draw(win)
			}
		}

		win.Update()
		<-fps
	}
}

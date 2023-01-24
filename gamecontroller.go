package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

func gameControl(win *pixelgl.Window) {

	setDebugWindow(win)

	backgroundPic := loadPicture("images/sprites/cenario01.png")
	backgroundSpr := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())
	shadePic := loadPicture("images/sprites/preto_60.png")
	shadeSpr := pixel.NewSprite(shadePic, shadePic.Bounds())

	caseTextBox := makeTextBox(
		"txt/cases.txt",
		pixel.V(win.Bounds().W()*(0.2), win.Bounds().H()*(1.0/30.0)),
		pixel.V(win.Bounds().W()*(0.8), win.Bounds().H()*(29.0/30.0)))
	propositionBoxes := makePropositionBoxes("txt/propositions.txt")

	botao := newButton("images/sprites/doc.png", 2)
	botao.setPosition(220, 50)

	var hoverButtons []*button
	hoverButtons = append(hoverButtons, &botao)

	fps := time.Tick(time.Second / 60)
	framesC := 0
	frames := 0
	second := time.Tick(time.Second)
	imd := imdraw.New(nil)
	imd.Color = colornames.Red

	// main loop
	var (
		mousePos           pixel.Vec
		screenCenterMatrix = pixel.IM.Moved(win.Bounds().Center())
		readingCaseTextBox = true
	)

	for !win.Closed() {

		mousePos = win.MousePosition()

		//events
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if readingCaseTextBox {
				if !caseTextBox.rect.Contains(mousePos) {
					readingCaseTextBox = false
				}
			} else {
				if botao.rect.Contains(mousePos) {
					readingCaseTextBox = true
					botao.state = 0
				} else {
					for i := range propositionBoxes.allBoxes {
						if propositionBoxes.allBoxes[i].txtBox.Bounds().Contains(mousePos) {
							propositionBoxes.allBoxes[i].click = !propositionBoxes.allBoxes[i].click
						}
					}
				}
			}
		}

		//draw
		win.Clear(colornames.Whitesmoke)
		backgroundSpr.Draw(win, screenCenterMatrix)

		for _, clickText := range propositionBoxes.allBoxes {
			clickText.txtBox.Draw(win, pixel.IM)
		}

		propositionBoxes.atkBorder.Draw(win)
		propositionBoxes.defBorder.Draw(win)
		botao.buttonSprs[botao.state].Draw(win, pixel.IM.Moved(botao.rect.Center()))

		if readingCaseTextBox {
			shadeSpr.Draw(win, screenCenterMatrix)
			caseTextBox.border.Draw(win)
			caseTextBox.bg.Draw(win)
			caseTextBox.txt.Draw(win, pixel.IM)
		}

		//hover
		if !readingCaseTextBox {
			for _, prop := range propositionBoxes.allBoxes {
				if prop.txtBox.Bounds().Contains(mousePos) || prop.click {
					imd.Clear()
					imd.Push(prop.txtBox.Bounds().Min, prop.txtBox.Bounds().Max)
					imd.Rectangle(0)
					imd.Draw(win)
				}
			}
			for _, button_ := range hoverButtons {
				if button_.rect.Contains(mousePos) {
					button_.state = 1
				} else {
					button_.state = 0
				}
			}
		}

		displayDebug("MousePos = "+mousePos.String(), 0)
		displayDebug(fmt.Sprintf("FPS = %d", frames), 1)
		displayDebug("readingCaseTextBox = "+strconv.FormatBool(readingCaseTextBox), 2)

		win.Update()

		//fps
		framesC++
		select {
		case <-second:
			frames = framesC
			framesC = 0
		default:
		}
		<-fps
	}
}

package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

func gameControl(win *pixelgl.Window) {

	setDebugWindow(win)

	backgroundPic := loadPicture("images/sprites/fundo.png")
	backgroundSpr := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())
	shadePic := loadPicture("images/sprites/preto_60.png")
	shadeSpr := pixel.NewSprite(shadePic, shadePic.Bounds())

	caseTextBox, guilty := makeTextBox(
		"txt/cases.txt",
		pixel.V(win.Bounds().W()*(0.2), win.Bounds().H()*(1.0/30.0)),
		pixel.V(win.Bounds().W()*(0.8), win.Bounds().H()*(29.0/30.0)))
	propositionBoxes := makePropositionBoxes("txt/propositions.txt")

	docButton := newButton("images/sprites/doc.png", 2)
	docButton.setPosition(60, 60)

	decisionButtons := makeDecisionButtons("images/sprites/culpado.png", "images/sprites/inocente.png", 2)
	decisionButtons.setPosition(win.Bounds().Center().X, win.Bounds().H()*0.395)
	gltyButton := decisionButtons.buttons[0]
	inncButton := decisionButtons.buttons[1]

	var hoverButtons []*button
	hoverButtons = append(hoverButtons, &docButton)

	fps := time.Tick(time.Second / 60)
	framesC := 0
	frames := 0
	second := time.Tick(time.Second)
	imd := imdraw.New(nil)
	imd.Color = colornames.Red

	// main loop
	var (
		mousePos             pixel.Vec
		screenCenterMatrix   = pixel.IM.Moved(win.Bounds().Center())
		backgroundMatrix     = screenCenterMatrix.Scaled(win.Bounds().Center(), win.Bounds().Max.X/backgroundPic.Bounds().Max.X)
		gltyButtonMatrix     = pixel.IM.Moved(gltyButton.rect.Center()).Scaled(gltyButton.rect.Center(), 1.5)
		inncButtonMatrix     = pixel.IM.Moved(inncButton.rect.Center()).Scaled(inncButton.rect.Center(), 1.5)
		readingCaseTextBox   = true
		endGame              = false
		acc                  int
		hitProps             float64
		total                float64
		endText              = text.New(win.Bounds().Center(), textAtlas70pt)
		decisionText         = text.New(win.Bounds().Center(), textAtlas70pt)
		winCondition         bool
		end                  bool
		decisionTextToCenter pixel.Vec
		endTextToCenter      pixel.Vec
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
				if docButton.rect.Contains(mousePos) {
					readingCaseTextBox = true
					docButton.state = 0
				} else {
					for i := range propositionBoxes.allBoxes {
						prop := &propositionBoxes.allBoxes[i]
						if prop.state.Bounds().Contains(mousePos) {
							prop.state = prop.txtStates[1]
							prop.click = !prop.click
						}
					}
					if !endGame {
						if decisionButtons.buttons[0].rect.Contains(mousePos) {
							end = true
							if guilty {
								winCondition = true
							} else {
								winCondition = false
							}
						} else if decisionButtons.buttons[1].rect.Contains(mousePos) {
							end = true
							if !guilty {
								winCondition = true
							} else {
								winCondition = false
							}
						}
						if end {
							if winCondition {
								decisionText.Color = colornames.Green
								decisionText.WriteString("GANHOU")
							} else {
								decisionText.Color = colornames.Red
								decisionText.WriteString("PERDEU")
							}
							endGame = true
							total = float64(len(propositionBoxes.allBoxes))
							for _, prop := range propositionBoxes.allBoxes {
								if (prop.trueProp && prop.click) || (!prop.trueProp && !prop.click) {
									hitProps++
								}
							}
							acc = int((hitProps / total) * 100)
							endText.Color = colornames.Black
							endText.WriteString(fmt.Sprintf("%v%% DE ACERTO", acc))
							decisionTextToCenter = pixel.V(-decisionText.Bounds().W()/2, (-decisionText.Bounds().H()/2)+(decisionText.LineHeight/2))
							endTextToCenter = pixel.V(-endText.Bounds().W()/2, (-endText.Bounds().H()/2)-(decisionText.LineHeight/2))
						}
					}
				}
			}
		}
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		//draw
		win.Clear(colornames.Whitesmoke)
		if !endGame {
			backgroundSpr.Draw(win, backgroundMatrix)

			for _, clickText := range propositionBoxes.allBoxes {
				clickText.state.Draw(win, pixel.IM)
			}

			gltyButton.buttonSprs[gltyButton.state].Draw(win, gltyButtonMatrix)
			inncButton.buttonSprs[inncButton.state].Draw(win, inncButtonMatrix)

			propositionBoxes.atkBorder.Draw(win)
			propositionBoxes.defBorder.Draw(win)
			docButton.buttonSprs[docButton.state].Draw(win, pixel.IM.Moved(docButton.rect.Center()))

			if readingCaseTextBox {
				shadeSpr.Draw(win, screenCenterMatrix)
				caseTextBox.border.Draw(win)
				caseTextBox.bg.Draw(win)
				caseTextBox.txt.Draw(win, pixel.IM)
			}

			//hover
			if !readingCaseTextBox {
				for i := range propositionBoxes.allBoxes {
					prop := &propositionBoxes.allBoxes[i]
					if !prop.click {
						if prop.state.Bounds().Contains(mousePos) {
							prop.state = prop.txtStates[1]
						} else {
							prop.state = prop.txtStates[0]
						}
					}
				}
				for _, button := range hoverButtons {
					if button.rect.Contains(mousePos) {
						button.state = 1
					} else {
						button.state = 0
					}
				}
				for _, button := range decisionButtons.buttons {
					if button.rect.Contains(mousePos) {
						button.state = 1
					} else {
						button.state = 0
					}
				}
			}
		} else {
			decisionText.Draw(win, pixel.IM.Moved(decisionTextToCenter))
			endText.Draw(win, pixel.IM.Moved(endTextToCenter))
		}

		displayDebug("MousePos = "+mousePos.String(), 0)
		displayDebug(fmt.Sprintf("FPS = %d", frames), 1)
		displayDebug("readingCaseTextBox = "+strconv.FormatBool(readingCaseTextBox), 2)
		displayDebug(fmt.Sprintf("1 proposition = %v", propositionBoxes.allBoxes[0].click), 3)
		displayDebug(fmt.Sprintf("2 proposition = %v", propositionBoxes.allBoxes[1].click), 4)
		displayDebug(fmt.Sprintf("3 proposition = %v", propositionBoxes.allBoxes[2].click), 5)
		displayDebug(fmt.Sprintf("4 proposition = %v", propositionBoxes.allBoxes[3].click), 6)
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

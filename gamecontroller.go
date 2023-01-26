package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

type endScreenText struct {
	uText         *text.Text
	bText         *text.Text
	uTextToCenter pixel.Matrix
	bTextToCenter pixel.Matrix
}

type char struct {
	spr    *pixel.Sprite
	rect   pixel.Rect
	matrix pixel.Matrix
}

func newChar(filePath string, posX float64, posY float64) char {
	pic := loadPicture(filePath)
	return char{
		spr:    pixel.NewSprite(pic, pic.Bounds()),
		rect:   pic.Bounds().Moved(pixel.V(posX, posY)),
		matrix: pixel.IM.Moved(pixel.V(posX, posY+(pic.Bounds().H()/2))),
	}
}

func endGame(w bool, t float64, b []clickableTxtBox) endScreenText {

	var (
		uText    = text.New(sCenter, textAtlas70pt)
		bText    = text.New(sCenter, textAtlas70pt)
		hitProps float64
		acc      int
	)

	if w {
		uText.Color = colornames.Green
		uText.WriteString("GANHOU")
	} else {
		uText.Color = colornames.Red
		uText.WriteString("PERDEU")
	}

	for _, prop := range b {
		if (prop.trueProp && prop.click) || (!prop.trueProp && !prop.click) {
			hitProps++
		}
	}

	acc = int((hitProps / t) * 100)
	bText.Color = colornames.Black
	bText.WriteString(fmt.Sprintf("%v%% DE ACERTO", acc))

	return endScreenText{
		uText:         uText,
		bText:         bText,
		uTextToCenter: pixel.IM.Moved(pixel.V(-uText.Bounds().W()/2, (-uText.Bounds().H()/2)+(uText.LineHeight/2))),
		bTextToCenter: pixel.IM.Moved(pixel.V(-bText.Bounds().W()/2, (-bText.Bounds().H()/2)-(uText.LineHeight/2))),
	}

}

func gameControl(win *pixelgl.Window) {

	setDebugWindow(win)

	backgroundPic := loadPicture("images/sprites/fundo.png")
	backgroundSpr := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())
	shadePic := loadPicture("images/sprites/preto_60.png")
	shadeSpr := pixel.NewSprite(shadePic, shadePic.Bounds())

	charJdg := newChar("images/sprites/judge.png", sWidth/2, 625)

	caseTextBox, guilty := makeTextBox(
		"txt/cases.txt",
		pixel.V(sWidth*(0.25), sHeight*(1.0/30.0)),
		pixel.V(sWidth*(0.75), sHeight*(29.0/30.0)),
	)

	propositionBoxes := makePropositionBoxes(
		"txt/propositions.txt",
		pixel.V(0, sHeight*(1.0/30.0)),
		pixel.V(sWidth*(0.26), sHeight*(29.0/30)),
	)

	docButton := newButton("images/sprites/doc.png", 2)
	docButton.setPosition(120, 60)

	arrowLButton := newButton("images/sprites/setaL.png", 2)
	arrowLButton.setPosition(0, (sHeight/2)-(arrowLButton.rect.H()/2))
	arrowRButton := newButton("images/sprites/setaR.png", 2)
	arrowRButton.setPosition(sWidth-arrowRButton.rect.W(), (sHeight/2)-(arrowRButton.rect.H()/2))

	decisionButtons := makeDecisionButtons("images/sprites/culpado.png", "images/sprites/inocente.png", 2)
	decisionButtons.setPosition(sWidth/2, sHeight*0.395)
	gltyButton := decisionButtons.buttons[0]
	inncButton := decisionButtons.buttons[1]

	var (
		hoverButtons       = []*button{&docButton}
		second             = time.Tick(time.Second)
		fps                = time.Tick(time.Second / 60)
		framesC            = 0
		frames             = 0
		mousePos           pixel.Vec
		screenCenterMatrix = pixel.IM.Moved(sCenter)
		backgroundMatrix   = screenCenterMatrix.Scaled(sCenter, sWidth/backgroundPic.Bounds().Max.X)
		gltyButtonMatrix   = pixel.IM.Moved(gltyButton.rect.Center()).Scaled(gltyButton.rect.Center(), 1.5)
		inncButtonMatrix   = pixel.IM.Moved(inncButton.rect.Center()).Scaled(inncButton.rect.Center(), 1.5)
		readingCaseTextBox = true
		endGameState       = false
		propTotal          = float64(len(propositionBoxes.allBoxes))
		endScreenTextBox   endScreenText
	)

	// main loop
	for !win.Closed() {

		mousePos = win.MousePosition()

		//events
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if readingCaseTextBox {
				if !caseTextBox.rect.Contains(mousePos) {
					readingCaseTextBox = false
				}
			} else {
				switch {
				case docButton.rect.Contains(mousePos):
					readingCaseTextBox = true
					docButton.state = 0
				case propositionBoxes.atkShow:
					for i := range propositionBoxes.atk {
						prop := &propositionBoxes.atk[i]
						if prop.state.Bounds().Contains(mousePos) {
							prop.state = prop.txtStates[1]
							prop.click = !prop.click
						}
					}
				case propositionBoxes.defShow:
					for i := range propositionBoxes.def {
						prop := &propositionBoxes.def[i]
						if prop.state.Bounds().Contains(mousePos) {
							prop.state = prop.txtStates[1]
							prop.click = !prop.click
						}
					}
				}
				if !endGameState {
					switch {
					case decisionButtons.buttons[0].rect.Contains(mousePos):
						endGameState = true
						if guilty {
							endScreenTextBox = endGame(true, propTotal, propositionBoxes.allBoxes)
						} else {
							endScreenTextBox = endGame(false, propTotal, propositionBoxes.allBoxes)
						}
					case decisionButtons.buttons[1].rect.Contains(mousePos):
						endGameState = true
						if !guilty {
							endScreenTextBox = endGame(true, propTotal, propositionBoxes.allBoxes)
						} else {
							endScreenTextBox = endGame(false, propTotal, propositionBoxes.allBoxes)
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
		if endGameState {
			endScreenTextBox.uText.Draw(win, endScreenTextBox.uTextToCenter)
			endScreenTextBox.bText.Draw(win, endScreenTextBox.bTextToCenter)
		} else {
			backgroundSpr.Draw(win, backgroundMatrix)

			charJdg.spr.Draw(win, charJdg.matrix)

			if propositionBoxes.atkShow {
				propositionBoxes.atkBg.Draw(win)
				propositionBoxes.atkBorder.Draw(win)
				for _, propText := range propositionBoxes.atk {
					propText.state.Draw(win, pixel.IM.Moved(propText.rect.Center()))
				}
			} else if propositionBoxes.defShow {
				propositionBoxes.defBg.Draw(win)
				propositionBoxes.defBorder.Draw(win)
				for _, propText := range propositionBoxes.def {
					propText.state.Draw(win, pixel.IM.Moved(propText.rect.Center()))
				}
			}

			gltyButton.buttonSprs[gltyButton.state].Draw(win, gltyButtonMatrix)
			inncButton.buttonSprs[inncButton.state].Draw(win, inncButtonMatrix)

			docButton.buttonSprs[docButton.state].Draw(win, pixel.IM.Moved(docButton.rect.Center()))

			arrowLButton.buttonSprs[arrowLButton.state].Draw(win, pixel.IM.Moved(arrowLButton.rect.Center()))
			arrowRButton.buttonSprs[arrowRButton.state].Draw(win, pixel.IM.Moved(arrowRButton.rect.Center()))

			if readingCaseTextBox {
				shadeSpr.Draw(win, screenCenterMatrix)
				caseTextBox.border.Draw(win)
				caseTextBox.bg.Draw(win)
				caseTextBox.txt.Draw(win, pixel.IM)
			}

			//hover
			if !readingCaseTextBox {
				switch {
				case !propositionBoxes.atkShow && arrowLButton.rect.Contains(mousePos):
					propositionBoxes.atkShow = true
					arrowLButton.rect = arrowLButton.rect.Moved(pixel.V(propositionBoxes.atkRect.W(), 0))
				case !propositionBoxes.defShow && arrowRButton.rect.Contains(mousePos):
					propositionBoxes.defShow = true
					arrowRButton.rect = arrowRButton.rect.Moved(pixel.V(-propositionBoxes.atkRect.W(), 0))
				case propositionBoxes.atkShow && !propositionBoxes.atkRect.Contains(mousePos):
					propositionBoxes.atkShow = false
					arrowLButton.rect = arrowLButton.rect.Moved(pixel.V(-propositionBoxes.atkRect.W(), 0))
				case propositionBoxes.defShow && !propositionBoxes.defRect.Contains(mousePos):
					propositionBoxes.defShow = false
					arrowRButton.rect = arrowRButton.rect.Moved(pixel.V(propositionBoxes.atkRect.W(), 0))
				}

				if propositionBoxes.atkShow {
					for i := range propositionBoxes.atk {
						prop := &propositionBoxes.atk[i]
						if !prop.click {
							if prop.state.Bounds().Contains(mousePos) {
								prop.state = prop.txtStates[1]
							} else {
								prop.state = prop.txtStates[0]
							}
						}
					}
				} else if propositionBoxes.defShow {
					for i := range propositionBoxes.def {
						prop := &propositionBoxes.def[i]
						if !prop.click {
							if prop.state.Bounds().Contains(mousePos) {
								prop.state = prop.txtStates[1]
							} else {
								prop.state = prop.txtStates[0]
							}
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
		}

		displayDebug("MousePos = "+mousePos.String(), 1)
		displayDebug(fmt.Sprintf("FPS = %d", frames), 2)
		displayDebug("readingCaseTextBox = "+strconv.FormatBool(readingCaseTextBox), 3)
		displayDebug(fmt.Sprintf("1 proposition = %v", propositionBoxes.atk[0].click), 4)
		displayDebug(fmt.Sprintf("2 proposition = %v", propositionBoxes.atk[1].click), 5)
		displayDebug(fmt.Sprintf("3 proposition = %v", propositionBoxes.atk[2].click), 6)
		displayDebug(fmt.Sprintf("4 proposition = %v", propositionBoxes.atk[3].click), 7)
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

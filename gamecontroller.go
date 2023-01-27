package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

type gameSetup struct {
	bg               pixel.Sprite
	bgMatrix         pixel.Matrix
	shade            pixel.Sprite
	shadeMatrix      pixel.Matrix
	chars            []char
	docButton        button
	docButtonMatrix  pixel.Matrix
	infoText         standardTxtBox
	propBoxes        propositionBoxes
	arrowLButton     button
	arrowRButton     button
	gltyButton       button
	gltyButtonMatrix pixel.Matrix
	inncButton       button
	inncButtonMatrix pixel.Matrix
	guilty           bool
}

func setup() gameSetup {
	var (
		pic pixel.Picture
		//w float64
		h float64
	)

	backgroundPic := loadPicture("images/sprites/fundo.png")
	backgroundSpr := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())
	shadePic := loadPicture("images/sprites/preto_60.png")
	shadeSpr := pixel.NewSprite(shadePic, shadePic.Bounds())

	chars := []char{
		newChar("images/sprites/judge.png", pixel.V(sWidth/2-(66), 625)),
		newChar("images/sprites/unknown.png", pixel.V(294, 550)),
		newChar("images/sprites/criminal.png", pixel.V(705, 150)),
		newChar("images/sprites/table.png", pixel.V((sWidth*0.25)-150, 40)),
		newChar("images/sprites/table.png", pixel.V((sWidth*0.75)-150, 40)),
		newChar("images/sprites/adv1.png", pixel.V((sWidth*0.25)-66, 120)),
		newChar("images/sprites/adv2.png", pixel.V((sWidth*0.75)-66, 120)),
	}

	docButton := newButton(loadPicture("images/sprites/docButton.png"), pixel.V(770, 650), 2)

	infoText, guilty := makeInfoText(
		"txt/cases.txt",
		pixel.V(sWidth*(0.25), sHeight*(1.0/30.0)),
		pixel.V(sWidth*(0.75), sHeight*(29.0/30.0)),
	)

	propBoxes := makePropositionBoxes(
		"txt/propositions.txt",
		pixel.V(0, sHeight*(1.0/30.0)),
		pixel.V(sWidth*(0.3), sHeight*(29.0/30)),
	)

	pic, _, h = loadPictureWH("images/sprites/setaL.png")
	arrowLButton := newButton(
		pic,
		pixel.V(0, (sHeight/2)-h/2),
		2,
	)
	pic, _, h = loadPictureWH("images/sprites/setaR.png")
	arrowRButton := newButton(
		pic,
		pixel.V(sWidth-88, (sHeight/2)-h/2),
		2,
	)

	decisionButtons := makeDecisionButtons(
		loadPicture("images/sprites/culpado.png"),
		loadPicture("images/sprites/inocente.png"),
		pixel.V(sWidth/2, sHeight*0.395),
		2,
	)

	return gameSetup{
		bg:               *backgroundSpr,
		bgMatrix:         pixel.IM.Moved(sCenter).Scaled(sCenter, sWidth/backgroundPic.Bounds().Max.X),
		shade:            *shadeSpr,
		shadeMatrix:      pixel.IM.Moved(sCenter).Scaled(sCenter, sWidth/shadePic.Bounds().Max.X),
		chars:            chars,
		docButton:        docButton,
		docButtonMatrix:  pixel.IM.Moved(docButton.rect.Center()),
		infoText:         infoText,
		propBoxes:        propBoxes,
		arrowLButton:     arrowLButton,
		arrowRButton:     arrowRButton,
		gltyButton:       *decisionButtons.buttons[0],
		gltyButtonMatrix: pixel.IM.Moved(decisionButtons.buttons[0].rect.Center()).Scaled(decisionButtons.buttons[0].rect.Center(), 1.5),
		inncButton:       *decisionButtons.buttons[1],
		inncButtonMatrix: pixel.IM.Moved(decisionButtons.buttons[1].rect.Center()).Scaled(decisionButtons.buttons[1].rect.Center(), 1.5),
		guilty:           guilty,
	}
}

func gameControl(win *pixelgl.Window) {

	setDebugWindow(win)
	g := setup()
	//catchEvents()
	//draw()
	//hover()
	//showDebugInfo()

	var (
		second             = time.Tick(time.Second)
		fps                = time.Tick(time.Second / 60)
		framesC            = 0
		frames             = 0
		mousePos           pixel.Vec
		readingCaseTextBox = true
		endGameState       = false
		endScreenTextBox   endScreenText
	)

	// main loop
	for !win.Closed() {

		mousePos = win.MousePosition()

		//events
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if readingCaseTextBox {
				if !g.infoText.rect.Contains(mousePos) {
					readingCaseTextBox = false
				}
			} else {
				switch {
				case g.docButton.rect.Contains(mousePos):
					readingCaseTextBox = true
					g.docButton.spr = &g.docButton.buttonSprs[0]
				case g.propBoxes.atkShow:
					for i := range g.propBoxes.atk {
						prop := g.propBoxes.atk[i]
						if (*prop.state).Bounds().Contains(mousePos) {
							prop.state = prop.txtStates[1]
							prop.click = !prop.click
						}
					}
				case g.propBoxes.defShow:
					for i := range g.propBoxes.def {
						prop := g.propBoxes.def[i]
						if (*prop.state).Bounds().Contains(mousePos) {
							prop.state = prop.txtStates[1]
							prop.click = !prop.click
						}
					}
				}
				if !endGameState {
					switch {
					case g.gltyButton.rect.Contains(mousePos):
						endGameState = true
						if g.guilty {
							endScreenTextBox = endGame(true, len(g.propBoxes.allBoxes), g.propBoxes.allBoxes)
						} else {
							endScreenTextBox = endGame(false, len(g.propBoxes.allBoxes), g.propBoxes.allBoxes)
						}
					case g.inncButton.rect.Contains(mousePos):
						endGameState = true
						if !g.guilty {
							endScreenTextBox = endGame(true, len(g.propBoxes.allBoxes), g.propBoxes.allBoxes)
						} else {
							endScreenTextBox = endGame(false, len(g.propBoxes.allBoxes), g.propBoxes.allBoxes)
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
			g.bg.Draw(win, g.bgMatrix)

			for _, char := range g.chars {
				char.spr.Draw(win, char.matrix)
			}

			(*g.docButton.spr).Draw(win, g.docButtonMatrix)

			(*g.arrowLButton.spr).Draw(win, pixel.IM.Moved(g.arrowLButton.rect.Center()))
			(*g.arrowRButton.spr).Draw(win, pixel.IM.Moved(g.arrowRButton.rect.Center()))

			if g.propBoxes.atkShow {
				g.propBoxes.atkBg.Draw(win)
				g.propBoxes.atkBorder.Draw(win)
				for _, propText := range g.propBoxes.atk {
					(*propText.state).Draw(win, pixel.IM.Moved(propText.rect.Center()))
				}
			} else if g.propBoxes.defShow {
				g.propBoxes.defBg.Draw(win)
				g.propBoxes.defBorder.Draw(win)
				for _, propText := range g.propBoxes.def {
					(*propText.state).Draw(win, pixel.IM.Moved(propText.rect.Center()))
				}
			}

			(*g.gltyButton.spr).Draw(win, g.gltyButtonMatrix)
			(*g.inncButton.spr).Draw(win, g.inncButtonMatrix)

			if readingCaseTextBox {
				g.shade.Draw(win, g.shadeMatrix)
				g.infoText.border.Draw(win)
				g.infoText.bg.Draw(win)
				g.infoText.txt.Draw(win, pixel.IM)
			}

			//hover
			if !readingCaseTextBox {
				switch {
				case !g.propBoxes.atkShow && g.arrowLButton.rect.Contains(mousePos):
					g.propBoxes.atkShow = true
					g.arrowLButton.spr = &g.arrowLButton.buttonSprs[1]
					g.arrowLButton.rect = g.arrowLButton.rect.Moved(pixel.V(g.propBoxes.atkRect.W(), 0))
				case !g.propBoxes.defShow && g.arrowRButton.rect.Contains(mousePos):
					g.propBoxes.defShow = true
					g.arrowRButton.spr = &g.arrowRButton.buttonSprs[1]
					g.arrowRButton.rect = g.arrowRButton.rect.Moved(pixel.V(-g.propBoxes.atkRect.W(), 0))
				case g.propBoxes.atkShow && !g.propBoxes.atkRect.Contains(mousePos):
					g.propBoxes.atkShow = false
					g.arrowLButton.spr = &g.arrowLButton.buttonSprs[0]
					g.arrowLButton.rect = g.arrowLButton.rect.Moved(pixel.V(-g.propBoxes.atkRect.W(), 0))
				case g.propBoxes.defShow && !g.propBoxes.defRect.Contains(mousePos):
					g.propBoxes.defShow = false
					g.arrowRButton.spr = &g.arrowRButton.buttonSprs[0]
					g.arrowRButton.rect = g.arrowRButton.rect.Moved(pixel.V(g.propBoxes.atkRect.W(), 0))
				}

				if g.propBoxes.atkShow {
					for i := range g.propBoxes.atk {
						prop := g.propBoxes.atk[i]
						if !prop.click {
							if prop.rect.Contains(mousePos) {
								prop.state = &prop.txtStates[1]
							} else {
								prop.state = &prop.txtStates[0]
							}
						}
					}
				} else if g.propBoxes.defShow {
					for i := range g.propBoxes.def {
						prop := &g.propBoxes.def[i]
						if !prop.click {
							if prop.rect.Contains(mousePos) {
								prop.state = &prop.txtStates[1]
							} else {
								prop.state = &prop.txtStates[0]
							}
						}
					}
				}

				if g.docButton.rect.Contains(mousePos) {
					g.docButton.spr = &g.docButton.buttonSprs[1]
				} else {
					g.docButton.spr = &g.docButton.buttonSprs[0]
				}

				if g.gltyButton.rect.Contains(mousePos) {
					g.gltyButton.spr = &g.gltyButton.buttonSprs[1]
				} else {
					g.gltyButton.spr = &g.gltyButton.buttonSprs[0]
				}
				if g.inncButton.rect.Contains(mousePos) {
					g.inncButton.spr = &g.inncButton.buttonSprs[1]
				} else {
					g.inncButton.spr = &g.inncButton.buttonSprs[0]
				}
			}
		}

		displayDebug("MousePos = "+mousePos.String(), 1)
		displayDebug(fmt.Sprintf("FPS = %d", frames), 2)
		displayDebug("readingCaseTextBox = "+strconv.FormatBool(readingCaseTextBox), 3)
		displayDebug(fmt.Sprintf("1 proposition = %v", g.propBoxes.atk[0].click), 4)
		displayDebug(fmt.Sprintf("2 proposition = %v", g.propBoxes.atk[1].click), 5)
		displayDebug(fmt.Sprintf("3 proposition = %v", g.propBoxes.atk[2].click), 6)
		displayDebug(fmt.Sprintf("4 proposition = %v", g.propBoxes.atk[3].click), 7)
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

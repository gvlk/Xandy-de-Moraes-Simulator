package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	face32pt      = loadTTF("fonts/upheavtt.ttf", 32)
	textAtlas32pt = text.NewAtlas(face32pt, text.ASCII, text.RangeTable(unicode.Latin))
	face38pt      = loadTTF("fonts/upheavtt.ttf", 38)
	textAtlas38pt = text.NewAtlas(face38pt, text.ASCII, text.RangeTable(unicode.Latin))
	face70pt      = loadTTF("fonts/upheavtt.ttf", 70)
	textAtlas70pt = text.NewAtlas(face70pt, text.ASCII, text.RangeTable(unicode.Latin))
)

const padding = 25

type standardTxtBox struct {
	txt    *text.Text
	border *imdraw.IMDraw
	bg     *imdraw.IMDraw
	rect   pixel.Rect
}

type propositionBoxes struct {
	atk       []clickableTxtBox
	def       []clickableTxtBox
	atkBorder *imdraw.IMDraw
	atkBg     *imdraw.IMDraw
	defBorder *imdraw.IMDraw
	defBg     *imdraw.IMDraw
	allBoxes  []clickableTxtBox
	atkRect   pixel.Rect
	defRect   pixel.Rect
	atkShow   bool
	defShow   bool
}

type clickableTxtBox struct {
	txtStates [2]text.Text
	state     **text.Text
	rect      pixel.Rect
	click     bool
	trueProp  bool
}

type endScreenText struct {
	uText         *text.Text
	bText         *text.Text
	uTextToCenter pixel.Matrix
	bTextToCenter pixel.Matrix
}

func makeInfoText(txtFileName string, min pixel.Vec, max pixel.Vec) (standardTxtBox, bool) {

	var (
		guilty = true
		ni     = false
	)

	textBorder := imdraw.New(nil)
	textBorder.Color = colornames.Darkgrey
	textBorder.Push(min, max)
	textBorder.Rectangle(5)

	textBg := imdraw.New(nil)
	textBg.Color = colornames.Whitesmoke
	textBg.Push(min.Add(pixel.V(2, 2)), max.Sub(pixel.V(2, 2)))
	textBg.Rectangle(0)

	data, err := os.ReadFile(txtFileName)
	if err != nil {
		panic(err)
	}
	words := strings.Fields(string(data))

	textBox := text.New(pixel.V(min.X+padding, max.Y-textAtlas38pt.LineHeight()), textAtlas38pt)
	textBox.Color = colornames.Black

	for i, word := range words {
		switch word {
		case ";nl":
			textBox.WriteString("\n")
			continue
		case ";ni":
			if !ni {
				ni = true
				textBox.WriteString("   - ")
			} else {
				ni = false
			}
			continue
		case ";c":
			continue
		case ";i":
			guilty = false
			continue
		default:
			if textBox.Dot.X+textBox.BoundsOf(words[i]).W() > max.X-padding {
				textBox.WriteString("\n")
				if ni {
					textBox.WriteString("      ")
				}
			}
			textBox.WriteString(word + " ")
		}

	}

	return standardTxtBox{
		txt:    textBox,
		border: textBorder,
		bg:     textBg,
		rect:   pixel.R(min.X, min.Y, max.X, max.Y),
	}, guilty
}

func makePropositionBox(propositions [][]string, min pixel.Vec, max pixel.Vec) ([]clickableTxtBox, *imdraw.IMDraw) {

	var (
		boxBorder   *imdraw.IMDraw
		mainTextBox *text.Text
		secuTextBox *text.Text
		textBoxes   []clickableTxtBox
		linesUsed   = 1.0
		trueProp    = false
	)

	boxBorder = imdraw.New(nil)
	boxBorder.Color = colornames.Black
	boxBorder.Push(min, max)
	boxBorder.Rectangle(5)

	for _, proposition := range propositions {
		mainTextBox = text.New(pixel.V(min.X+10, max.Y-(linesUsed*textAtlas32pt.LineHeight())), textAtlas32pt)
		secuTextBox = text.New(pixel.V(min.X+10, max.Y-(linesUsed*textAtlas32pt.LineHeight())), textAtlas32pt)
		mainTextBox.Color = colornames.Black
		secuTextBox.Color = colornames.Orangered
		for j, word := range proposition {
			if word == ";t" {
				trueProp = true
				continue
			}
			if mainTextBox.Dot.X+mainTextBox.BoundsOf(proposition[j]).W() > max.X-padding {
				mainTextBox.WriteString("\n")
				secuTextBox.WriteString("\n")
				linesUsed++
			}
			mainTextBox.WriteString(word + " ")
			secuTextBox.WriteString(word + " ")
		}
		linesUsed += 0.8
		textBoxes = append(textBoxes, clickableTxtBox{
			txtStates: [2]text.Text{*mainTextBox, *secuTextBox},
			state:     &mainTextBox,
			click:     false,
			trueProp:  trueProp,
		})
		linesUsed++
		trueProp = false
	}
	return textBoxes, boxBorder
}

func makePropositionBoxes(txtFileName string, min pixel.Vec, max pixel.Vec) propositionBoxes {
	data, err := os.ReadFile(txtFileName)
	if err != nil {
		panic(err)
	}

	var (
		words    = strings.Fields(string(data))
		atk      [][]string
		def      [][]string
		atkwords []string
		defwords []string
		side     = false
		allBoxes []clickableTxtBox
	)

	for _, word := range words {
		switch word {
		case ";p":
			if !side {
				atk = append(atk, atkwords)
				atkwords = nil
			} else {
				def = append(def, defwords)
				defwords = nil
			}
		case ";def":
			side = true
		default:
			if !side {
				atkwords = append(atkwords, word)
			} else {
				defwords = append(defwords, word)
			}
		}
	}

	atkBoxes, atkBorder := makePropositionBox(atk, min, max)
	defBoxes, defBorder := makePropositionBox(def, pixel.V(sWidth-max.X, min.Y), pixel.V(sWidth, max.Y))
	allBoxes = append(allBoxes, atkBoxes...)
	allBoxes = append(allBoxes, defBoxes...)

	atkBg := imdraw.New(nil)
	atkBg.Color = pixel.RGBA{R: 0.97, G: 0.97, B: 0.97, A: 0.97}
	atkBg.Push(min.Add(pixel.V(2, 2)), max.Sub(pixel.V(2, 2)))
	atkBg.Rectangle(0)

	defBg := imdraw.New(nil)
	defBg.Color = pixel.RGBA{R: 0.97, G: 0.97, B: 0.97, A: 0.97}
	defBg.Push(pixel.V(sWidth-max.X, min.Y).Add(pixel.V(2, 2)), pixel.V(sWidth, max.Y).Sub(pixel.V(2, 2)))
	defBg.Rectangle(0)

	return propositionBoxes{
		atk:       atkBoxes,
		def:       defBoxes,
		atkBorder: atkBorder,
		atkBg:     atkBg,
		defBorder: defBorder,
		defBg:     defBg,
		allBoxes:  allBoxes,
		atkRect:   pixel.R(min.X, min.Y, max.X, max.Y),
		defRect:   pixel.R(sWidth-max.X, min.Y, sWidth, max.Y),
		atkShow:   false,
		defShow:   false,
	}
}

func endGame(w bool, t int, b []clickableTxtBox) endScreenText {

	var (
		uText    = text.New(sCenter, textAtlas70pt)
		bText    = text.New(sCenter, textAtlas70pt)
		hitProps int
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

	acc = int((float64(hitProps) / float64(t)) * 100)
	bText.Color = colornames.Black
	bText.WriteString(fmt.Sprintf("%v%% DE ACERTO", acc))

	return endScreenText{
		uText:         uText,
		bText:         bText,
		uTextToCenter: pixel.IM.Moved(pixel.V(-uText.Bounds().W()/2, (-uText.Bounds().H()/2)+(uText.LineHeight/2))),
		bTextToCenter: pixel.IM.Moved(pixel.V(-bText.Bounds().W()/2, (-bText.Bounds().H()/2)-(uText.LineHeight/2))),
	}

}

func loadTTF(path string, size float64) font.Face {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	font_, err := truetype.Parse(bytes)
	if err != nil {
		panic(err)
	}

	return truetype.NewFace(font_, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

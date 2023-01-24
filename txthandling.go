package main

import (
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
	face38pt      = loadTTF("fonts/upheavtt.ttf", 38)
	textAtlas38pt = text.NewAtlas(face38pt, text.ASCII, text.RangeTable(unicode.Latin))
	face24pt      = loadTTF("fonts/upheavtt.ttf", 24)
	textAtlas24pt = text.NewAtlas(face24pt, text.ASCII, text.RangeTable(unicode.Latin))
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
	defBorder *imdraw.IMDraw
	allBoxes  []clickableTxtBox
}

type clickableTxtBox struct {
	txtBox *text.Text
	click  bool
}

func makeTextBox(txtFileName string, min pixel.Vec, max pixel.Vec) standardTxtBox {

	textBorder := imdraw.New(nil)
	textBorder.Color = colornames.Darkgrey
	textBorder.Push(min, max)
	textBorder.Rectangle(5)

	textBG := imdraw.New(nil)
	textBG.Color = colornames.Whitesmoke
	textBG.Push(min.Add(pixel.V(2, 2)), max.Sub(pixel.V(2, 2)))
	textBG.Rectangle(0)

	data, err := os.ReadFile(txtFileName)
	if err != nil {
		panic(err)
	}
	words := strings.Fields(string(data))

	textBox := text.New(pixel.V(min.X+padding, max.Y-textAtlas38pt.LineHeight()), textAtlas38pt)
	textBox.Color = colornames.Black

	ni := false
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
		bg:     textBG,
		rect:   pixel.R(min.X, min.Y, max.X, max.Y),
	}
}

func makePropositionBox(propositions [][]string, min pixel.Vec, max pixel.Vec) ([]clickableTxtBox, *imdraw.IMDraw) {

	var (
		boxBorder *imdraw.IMDraw
		textBoxes []clickableTxtBox
		linesUsed = 1.0
	)

	boxBorder = imdraw.New(nil)
	boxBorder.Color = colornames.Black
	boxBorder.Push(min, max)
	boxBorder.Rectangle(5)

	for _, proposition := range propositions {
		textBox := text.New(pixel.V(min.X+10, max.Y-(linesUsed*textAtlas24pt.LineHeight())), textAtlas24pt)
		textBox.Color = colornames.Black
		for j, word := range proposition {
			if textBox.Dot.X+textBox.BoundsOf(proposition[j]).W() > max.X-padding {
				textBox.WriteString("\n")
				linesUsed++
			}
			textBox.WriteString(word + " ")
		}
		textBoxes = append(textBoxes, clickableTxtBox{txtBox: textBox, click: false})
		linesUsed++
	}
	return textBoxes, boxBorder
}

func makePropositionBoxes(txtFileName string) propositionBoxes {
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

	atkBoxes, atkBorder := makePropositionBox(atk, pixel.V(40, 400), pixel.V(580, 670))
	defBoxes, defBorder := makePropositionBox(def, pixel.V(895, 400), pixel.V(1435, 670))
	allBoxes = append(allBoxes, atkBoxes...)
	allBoxes = append(allBoxes, defBoxes...)

	return propositionBoxes{
		atk:       atkBoxes,
		def:       defBoxes,
		atkBorder: atkBorder,
		defBorder: defBorder,
		allBoxes:  allBoxes,
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

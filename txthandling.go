package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image/color"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	face      = loadTTF("fonts/upheavtt.ttf", 20)
	textAtlas = text.NewAtlas(face, text.ASCII, text.RangeTable(unicode.Latin))
)

func makeTextBox(txtFileName string, leftBottomVertice pixel.Vec, RightUpperVertice pixel.Vec) (*text.Text, *imdraw.IMDraw) {

	txtRectangle := imdraw.New(nil)
	txtRectangle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	txtRectangle.Push(leftBottomVertice, RightUpperVertice)
	txtRectangle.Rectangle(5)

	data, err := os.ReadFile(txtFileName)
	if err != nil {
		panic(err)
	}
	words := strings.Fields(string(data))

	txtBox := text.New(pixel.V(leftBottomVertice.X+10, RightUpperVertice.Y-textAtlas.LineHeight()), textAtlas)
	txtBox.Color = colornames.Black

	for i, word := range words {
		if txtBox.Dot.X+txtBox.BoundsOf(words[i]).W() > RightUpperVertice.X-10 {
			txtBox.WriteString("\n")
		}
		txtBox.WriteString(word + " ")
	}

	return txtBox, txtRectangle
}

func makeMultiTextBox(txtFileName string, leftBottomVertice pixel.Vec, RightUpperVertice pixel.Vec) ([]*text.Text, *imdraw.IMDraw) {

	txtRectangle := imdraw.New(nil)
	txtRectangle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	txtRectangle.Push(leftBottomVertice, RightUpperVertice)
	txtRectangle.Rectangle(5)

	data, err := os.ReadFile(txtFileName)
	if err != nil {
		panic(err)
	}
	words := strings.Fields(string(data))

	var txtBoxes []*text.Text

	for i, word := range words {
		txtBox := text.New(pixel.V(
			leftBottomVertice.X+10,
			RightUpperVertice.Y-(float64(i+1)*textAtlas.LineHeight())), textAtlas)
		txtBox.Color = colornames.Black
		txtBox.WriteString(word)
		txtBoxes = append(txtBoxes, txtBox)
	}

	return txtBoxes, txtRectangle
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

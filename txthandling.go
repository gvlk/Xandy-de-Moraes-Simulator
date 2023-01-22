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
	//	fragmentShader = `
	//#version 330 core
	//in vec2  vTexCoords;
	//out vec4 fragColor;
	//uniform vec4 uTexBounds;
	//uniform sampler2D uTexture;
	//void main() {
	//	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	//	// Sum our 3 color channels
	//	float sum  = texture(uTexture, t).r;
	//	      sum += texture(uTexture, t).g;
	//	      sum += texture(uTexture, t).b;
	//	// Divide by 3, and set the output to the result
	//	vec4 color = vec4( 0.1, 0.2, 0.5, 1.0);
	//	fragColor = color;
	//}
	//`
)

type clickableTxtBox struct {
	txtBox *text.Text
	click  bool
}

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

func makeClickable(txtBoxes []*text.Text) []clickableTxtBox {

	var clickableTxtBoxes []clickableTxtBox

	for _, txtBox := range txtBoxes {
		x := clickableTxtBox{txtBox: txtBox, click: false}
		clickableTxtBoxes = append(clickableTxtBoxes, x)
	}

	return clickableTxtBoxes
}

func makeMultiClickTextBox(txtFileName string, leftBottomVertice pixel.Vec, RightUpperVertice pixel.Vec) ([]clickableTxtBox, *imdraw.IMDraw) {
	txtBoxes, txtRectangle := makeMultiTextBox(txtFileName, leftBottomVertice, RightUpperVertice)
	return makeClickable(txtBoxes), txtRectangle
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

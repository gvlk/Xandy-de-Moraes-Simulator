package main

import "github.com/faiface/pixel"

type button struct {
	buttonSprs   []*pixel.Sprite
	buttonFrames []pixel.Rect
	rect         pixel.Rect
	state        int
}

func newButton(imgPath string, sprNum float64) button {

	var (
		spritesheet  = loadPicture(imgPath)
		buttonSprs   []*pixel.Sprite
		buttonFrames []pixel.Rect
		sprLength    = spritesheet.Bounds().W() / sprNum
		sprHeight    = spritesheet.Bounds().H()
	)

	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += sprLength {
		buttonFrames = append(buttonFrames, pixel.R(x, 0, x+sprLength, sprHeight))
	}
	for x := 0; x < int(sprNum); x++ {
		buttonSpr := pixel.NewSprite(spritesheet, buttonFrames[x])
		buttonSprs = append(buttonSprs, buttonSpr)
	}

	return button{
		buttonSprs:   buttonSprs,
		buttonFrames: buttonFrames,
		rect:         buttonFrames[0],
		state:        0,
	}
}

func (b *button) setPosition(x float64, y float64) {
	b.rect = b.rect.Moved(pixel.V(x, y))
}

package main

import "github.com/faiface/pixel"

type decisionButtons struct {
	buttons [2]*button
}

type button struct {
	buttonSprs []*pixel.Sprite
	spr        **pixel.Sprite
	rect       pixel.Rect
}

func newButton(pic pixel.Picture, min pixel.Vec, sprNum float64) button {

	var (
		buttonSprs   []*pixel.Sprite
		buttonFrames []pixel.Rect
		sprLength    = pic.Bounds().W() / sprNum
		sprHeight    = pic.Bounds().H()
	)

	for x := pic.Bounds().Min.X; x < pic.Bounds().Max.X; x += sprLength {
		buttonFrames = append(buttonFrames, pixel.R(x, 0, x+sprLength, sprHeight))
	}
	for x := 0; x < int(sprNum); x++ {
		buttonSpr := pixel.NewSprite(pic, buttonFrames[x])
		buttonSprs = append(buttonSprs, buttonSpr)
	}

	return button{
		buttonSprs: buttonSprs,
		spr:        &buttonSprs[0],
		rect:       buttonFrames[0].Moved(min),
	}
}

func makeDecisionButtons(gltyPic pixel.Picture, inncPic pixel.Picture, center pixel.Vec, sprNum float64) decisionButtons {
	gltyButton := newButton(gltyPic, pixel.V(center.X-gltyPic.Bounds().W()-45, center.Y-(gltyPic.Bounds().H()/2)), sprNum)
	inncButton := newButton(inncPic, pixel.V(center.X+45, center.Y-(inncPic.Bounds().H()/2)), sprNum)

	return decisionButtons{
		buttons: [2]*button{&gltyButton, &inncButton},
	}
}

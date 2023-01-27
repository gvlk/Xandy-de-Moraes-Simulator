package main

import (
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"os"
)

type char struct {
	spr    *pixel.Sprite
	rect   pixel.Rect
	matrix pixel.Matrix
}

func loadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img)
}

func loadPictureWH(path string) (pixel.Picture, float64, float64) {
	pic := loadPicture(path)
	return pic, pic.Bounds().W(), pic.Bounds().H()
}

func newChar(path string, min pixel.Vec) char {
	pic := loadPicture(path)
	return char{
		spr:    pixel.NewSprite(pic, pic.Bounds()),
		rect:   pic.Bounds().Moved(min),
		matrix: pixel.IM.Moved(pixel.V(min.X+(pic.Bounds().W()/2), min.Y+(pic.Bounds().H()/2))),
	}
}

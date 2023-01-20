package main

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io"
	"os"
)

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font_, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font_, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

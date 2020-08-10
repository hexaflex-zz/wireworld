// Package ui defines some rudimentary UI components.
package ui

import (
	"image"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	fontRegularSize = 10
	fontRegularDPI  = 96
)

var (
	fontRegular           *freetype.Context
	fontRegularLineHeight int
)

func init() {
	ttf, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}

	fontRegular = freetype.NewContext()
	fontRegular.SetDPI(fontRegularDPI)
	fontRegular.SetFont(ttf)
	fontRegular.SetFontSize(fontRegularSize)
	fontRegular.SetHinting(font.HintingFull)
	fontRegular.SetSrc(image.Black)

	fontRegularLineHeight = int(fontRegular.PointToFixed(fontRegularSize+2) >> 6)
}

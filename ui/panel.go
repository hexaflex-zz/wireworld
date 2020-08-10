package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"wireworld/resources"
	"wireworld/util"

	"github.com/golang/freetype"
)

// InfoPanel defines a rectangular panel with debug information.
type InfoPanel struct {
	x, y, w, h     int
	image          *image.RGBA
	textureChanged bool
}

// NewPanel creates a new debug info panel.
func NewInfoPanel() *InfoPanel {
	var p InfoPanel
	return &p
}

func (p *InfoPanel) Release() {
	p.image = nil
}

// Resize resizes and positions the button.
func (p *InfoPanel) Resize(x, y, w, h int) {
	if p.w != w || p.h != h {
		p.textureChanged = true
		p.image = image.NewRGBA(image.Rect(0, 0, util.Pow2(w), util.Pow2(h)))
		p.w = w
		p.h = h
	}

	p.x = x
	p.y = y
}

// Clear clears all panel contents.
func (p *InfoPanel) Clear() {
	clr := image.NewUniform(color.RGBA{0xb3, 0xb3, 0xb3, 0xff})
	draw.Draw(p.image, p.image.Bounds(), clr, image.ZP, draw.Src)
	p.textureChanged = true
}

// Print sets the given line to the specified formatted content.
func (p *InfoPanel) Print(line int, v string, argv ...interface{}) {
	v = fmt.Sprintf(v, argv...)
	if len(v) == 0 {
		return
	}

	x := p.x + 5
	y := p.y + (line+1)*fontRegularLineHeight

	fontRegular.SetClip(p.image.Bounds())
	fontRegular.SetDst(p.image)
	fontRegular.DrawString(v, freetype.Pt(x, y))

	p.textureChanged = true
}

func (p *InfoPanel) Draw(mp *util.Mat4) {
	mvp := mp.Copy()
	mvp.Mul(util.Mat4Translate(float32(p.x), float32(p.y), 0))
	mvp.Mul(util.Mat4Scale(float32(p.w), float32(p.h), 0))

	s := resources.GetShader("Panel")
	s.Use()
	s.SetMat16("mvp", mvp[:])

	m := resources.GetMesh("Panel").(*resources.TexturedQuadMesh)

	// Upload texture, if applicable.
	if p.textureChanged {
		p.textureChanged = false
		m.CommitTexture(p.image, p.w, p.h)
	}

	m.Draw()
}

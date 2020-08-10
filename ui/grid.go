package ui

import (
	"wireworld/resources"
	"wireworld/util"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Grid extends a canvas with grid rendering and snapping functionality.
type Grid struct {
	*Canvas

	uniformInvalid bool
	gridInvalid    bool
	gridVisible    bool
}

// NewGrid creates a new grid ontop of a canvas
func NewGrid() *Grid {
	return &Grid{
		Canvas:         NewCanvas(),
		uniformInvalid: true,
		gridInvalid:    true,
		gridVisible:    true,
	}
}

func (g *Grid) Resize(w, h int) {
	g.Canvas.Resize(w, h)
	g.uniformInvalid = true
	g.gridInvalid = true
}

func (g *Grid) Scroll(x, y float64) {
	g.Canvas.Scroll(x, y)
	g.uniformInvalid = true
	g.gridInvalid = true
}

func (g *Grid) SetPanning(v bool) {
	g.Canvas.SetPanning(v)
	g.uniformInvalid = true
}

// ToggleGridVisible toggles visibility of the grid background.
// Returns the new state.
func (g *Grid) ToggleGridVisible() bool {
	g.gridVisible = !g.gridVisible
	return g.gridVisible
}

func (g *Grid) Draw(mp *util.Mat4) {
	g.Canvas.Draw(mp)

	if !g.gridVisible {
		return
	}

	s := resources.GetShader("Grid")
	s.Use()

	// Upload shader uniforms if needed.
	if g.uniformInvalid {
		g.uniformInvalid = false

		cs := g.Zoom()
		px := float32(g.origin[0] % cs)
		py := float32(g.origin[1] % cs)

		mvp := mp.Copy()
		mvp.Mul(util.Mat4Translate(px, py, 0))

		s.SetMat16("mvp", (*mvp)[:])
	}

	m := resources.GetMesh("Grid")

	// Recreate the grid mesh if needed.
	if g.gridInvalid {
		g.gridInvalid = false
		g.createGrid(m)
	}

	m.Draw()
}

// createGrid regenerates and uploads the grid mesh, based on the
// current zoom factor and viewport dimensions. It consists of a set
// of horizontal and vertical lines spanning the full width or height
// of the viewport. Each separated by the zoomed cell size.
func (g *Grid) createGrid(m resources.Mesh) {
	cs := g.Zoom()
	w, h := g.Viewport()
	w += cs * 2
	h += cs * 2

	lines := make([]float32, 0, int((w/cs)+(h/cs))*4)

	// Horizontal lines.
	for y := 0; y < h; y += cs {
		lines = append(lines,
			float32(-cs),
			float32(y),
			float32(w+cs),
			float32(y))
	}

	// Vertical lines.
	for x := 0; x < w; x += cs {
		lines = append(lines,
			float32(x),
			float32(-cs),
			float32(x),
			float32(h+cs))
	}

	m.Commitfv(lines, gl.STREAM_DRAW)
}

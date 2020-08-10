package ui

import (
	"math"

	"wireworld/resources"
	"wireworld/sim"
	"wireworld/util"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	ZoomMin     = 1
	ZoomMax     = 30
	ZoomDefault = 15
)

// Canvas facilitates panning and zooming and tracks mouse input.
type Canvas struct {
	mousePosition [2]int
	mouseDelta    [2]int
	origin        [2]int
	viewport      [2]int
	zoom          int
	panning       bool
}

// NewCanvas creates a new canvas
func NewCanvas() *Canvas {
	return &Canvas{
		mousePosition: [2]int{0, 0},
		mouseDelta:    [2]int{0, 0},
		origin:        [2]int{0, 0},
		viewport:      [2]int{1, 1},
		zoom:          ZoomDefault,
		panning:       false,
	}
}

func (c *Canvas) Draw(mp *util.Mat4) {
	s := resources.GetShader("CellRenderer")
	s.Use()
	s.Set1f("alpha", 1.0)
	c.setCellMVP(s, mp, c.origin[0], c.origin[1])

	m := resources.GetMesh("CellRenderer")

	// Upload cell data if it has changed.
	if sim.CellsChanged() {
		m.Commitiv(sim.Cells(), gl.STREAM_DRAW)
	}

	m.Draw()
}

// setCellMVP computes the cell MVP matrix for the given shader
// and position.
func (c *Canvas) setCellMVP(s *resources.Shader, mp *util.Mat4, x, y int) {
	w, h := c.viewport[0], c.viewport[1]
	z := float32(c.zoom)

	mvp := mp.Copy()
	mvp.Mul(util.Mat4Translate(float32(x), float32(y), 0))
	mvp.Mul(util.Mat4Scale(z, z, 0))

	s.SetMat16("mvp", (*mvp)[:])
	s.Set2f("cellSize", z/float32(w)*2, z/float32(h)*2)
}

// MouseDelta returns the current mouse delta.
func (c *Canvas) MouseDelta() (int, int) {
	return c.mouseDelta[0], c.mouseDelta[1]
}

// MousePosition returns the current mouse position.
func (c *Canvas) MousePosition() (int, int) {
	return c.mousePosition[0], c.mousePosition[1]
}

// Viewport returns the viewport width and height.
func (c *Canvas) Viewport() (int, int) {
	return c.viewport[0], c.viewport[1]
}

// SetPanning sets the panning flag. Meaning we are either dragging
// the scene around or not.
func (c *Canvas) SetPanning(v bool) {
	c.panning = v
}

// ScrollTo scrolls the viewport to the given, absolute position.
func (c *Canvas) ScrollTo(x, y int) {
	c.origin[0] = x
	c.origin[1] = y
}

// Zoom returns the current zoom factor.
func (c *Canvas) Zoom() int {
	return c.zoom
}

// SetZoom sets the current zoom factor.
func (c *Canvas) SetZoom(v int) {
	c.zoom = util.Clampi(v, ZoomMin, ZoomMax)
}

func (c *Canvas) Resize(w, h int) {
	c.viewport[0] = w
	c.viewport[1] = h
}

func (c *Canvas) Scroll(x, y float64) {
	ox, oy := float32(c.origin[0]), float32(c.origin[1])
	zf := float32(c.zoom)
	zd := float32(y / math.Abs(y))
	fx := float32(c.mousePosition[0]) - ox
	fy := float32(c.mousePosition[1]) - oy

	oldfx := fx / zf
	oldfy := fy / zf

	c.zoom = util.Clampi(int(zf+zd), ZoomMin, ZoomMax)
	zf = float32(c.zoom)

	c.origin[0] = int(fx - ((oldfx * zf) - ox))
	c.origin[1] = int(fy - ((oldfy * zf) - oy))
}

func (c *Canvas) MouseMove(x, y float64) {
	c.mouseDelta[0] = c.mousePosition[0] - int(x)
	c.mouseDelta[1] = c.mousePosition[1] - int(y)
	c.mousePosition[0] = int(x)
	c.mousePosition[1] = int(y)

	if c.panning {
		c.origin[0] -= c.mouseDelta[0]
		c.origin[1] -= c.mouseDelta[1]
	}
}

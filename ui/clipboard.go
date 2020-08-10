package ui

import (
	"wireworld/resources"
	"wireworld/sim"
	"wireworld/util"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Clipboard extends a cell selector by providing clipboard functionality.
// It implements cut, copy and paste.
type Clipboard struct {
	*CellSelector

	clipboard        sim.CellList
	uniformChanged   bool // Shader uniforms need to be updated?
	clipboardChanged bool // Contents need to be recomitted to GPU?
	drawClipboard    bool // Draw clipboard contents?
}

// NewClipboard creates a new clipboard ontop of a cell selector.
func NewClipboard() *Clipboard {
	return &Clipboard{
		CellSelector:     NewCellSelector(),
		clipboardChanged: false,
		uniformChanged:   false,
		drawClipboard:    true,
	}
}

func (c *Clipboard) MouseMove(x, y float64) {
	c.CellSelector.MouseMove(x, y)
	c.uniformChanged = true
}

func (c *Clipboard) Scroll(x, y float64) {
	c.CellSelector.Scroll(x, y)
	c.uniformChanged = true
}

func (c *Clipboard) SetPanning(v bool) {
	c.CellSelector.SetPanning(v)
	c.uniformChanged = true
}

// ToggleDrawClipboard toggles drawing of the clipboard contents.
func (c *Clipboard) ToggleDrawClipboard() {
	c.drawClipboard = !c.drawClipboard
}

func (c *Clipboard) Draw(mp *util.Mat4) {
	c.CellSelector.Draw(mp)

	if !c.drawClipboard || c.clipboard.Len() == 0 {
		return
	}

	s := resources.GetShader("CellRenderer")
	s.Use()
	s.Set1f("alpha", 0.5)

	if c.uniformChanged {
		z := c.Zoom() / 2
		c.setCellMVP(s, mp,
			c.mousePosition[0]-z,
			c.mousePosition[1]-z)
		c.uniformChanged = false
	}

	m := resources.GetMesh("Clipboard")

	// Recommit cells to GPU if needed.
	if c.clipboardChanged {
		c.clipboardChanged = false
		m.Commitiv(c.clipboard, gl.STREAM_DRAW)
	}

	m.Draw()
}

// ClipboardClear empties the clipboard.
func (c *Clipboard) ClipboardClear() {
	c.clipboard = nil
}

// ClipboardCut copies the current selection to the clipboard and
// then deletes the selected cells from the simulation.
func (c *Clipboard) ClipboardCut() {
	c.ClipboardCopy()
	c.SelectionDelete()
	c.SelectionClear()
}

// ClipboardCopy copies the current cell selection to the clipboard.
func (c *Clipboard) ClipboardCopy() {
	c.clipboard = make(sim.CellList, len(c.selection))
	copy(c.clipboard, c.selection)

	// Treat the selection as a rectangle. Find the smallest
	// X and Y coordinate values.
	minx := int32(1<<31 - 1)
	miny := int32(1<<31 - 1)

	for i := 0; i < len(c.clipboard)-2; i += 3 {
		if c.clipboard[i+0] < minx {
			minx = c.clipboard[i+0]
		}

		if c.clipboard[i+1] < miny {
			miny = c.clipboard[i+1]
		}
	}

	// Offset the cell positions to be located at 0/0, while
	// preserving the relative distance between each cell.
	for i := 0; i < len(c.clipboard)-2; i += 3 {
		c.clipboard[i+0] -= minx
		c.clipboard[i+1] -= miny
	}

	c.clipboardChanged = true
}

// ClipboardPaste pastes the current cell selection from the clipboard,
// to the simulation at the current cursor position.
func (c *Clipboard) ClipboardPaste() {
	if c.clipboard.Len() == 0 {
		return
	}

	x, y := c.HoverTarget()
	sim.Load(x, y, c.clipboard)
}

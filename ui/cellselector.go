package ui

import (
	"image"
	"math"

	"wireworld/resources"
	"wireworld/sim"
	"wireworld/util"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// CellSelector facilitates the selection of one or more cells.
// It maintains the current selection and provides facilities to
// draw the selection and the selection rectangle, if it is being drawn.
// It builds on the CellRenderer functionality.
type CellSelector struct {
	*Grid

	selectionStart   *image.Point // beginning of selection rect.
	selection        sim.CellList // Current selection.
	selecting        bool         // Are we drawing a selection rect?
	addSelection     bool         // Are we adding to existing selection?
	selectionChanged bool         // Selection has changed and should be re-comitted o GPU.
}

// NewCellSelector creates a new cell selector ontop of a cell renderer.
func NewCellSelector() *CellSelector {
	return &CellSelector{
		Grid:             NewGrid(),
		selectionChanged: false,
		selecting:        false,
		selectionStart:   nil,
	}
}

func (c *CellSelector) Scroll(x, y float64) {
	c.Grid.Scroll(x, y)
	c.uniformInvalid = true
}

func (c *CellSelector) Draw(mp *util.Mat4) {
	c.Grid.Draw(mp)

	c.drawRectangle(mp)
	c.drawSelection(mp)
}

// drawRectangle draws the selection rectangle.
func (c *CellSelector) drawRectangle(mp *util.Mat4) {
	if !c.selecting || c.selectionStart == nil {
		return
	}

	x := c.selectionStart.X
	y := c.selectionStart.Y
	w := c.mousePosition[0] - x
	h := c.mousePosition[1] - y

	mvp := mp.Copy()
	mvp.Mul(util.Mat4Translate(float32(x), float32(y), 0))
	mvp.Mul(util.Mat4Scale(float32(w), float32(h), 0))

	s := resources.GetShader("CellSelectorRect")
	s.Use()
	s.SetMat16("mvp", (*mvp)[:])

	m := resources.GetMesh("CellSelectorRect")
	m.Draw()
}

func (c *CellSelector) drawSelection(mp *util.Mat4) {
	if c.selection.Len() == 0 {
		return
	}

	s := resources.GetShader("CellSelectorCells")
	s.Use()
	c.setCellMVP(s, mp, c.origin[0], c.origin[1])

	m := resources.GetMesh("CellSelectorCells")

	// Upload cell data if it has changed.
	if c.selectionChanged {
		c.selectionChanged = false
		m.Commitiv(c.selection, gl.STREAM_DRAW)
	}

	m.Draw()
}

// HoverTarget returns the cell coordinates at the current
// mouse cursor position.
func (c *CellSelector) HoverTarget() (int32, int32) {
	z := float64(c.Zoom())
	x := math.Floor(float64(c.mousePosition[0]-c.origin[0]) / z)
	y := math.Floor(float64(c.mousePosition[1]-c.origin[1]) / z)
	return int32(x), int32(y)
}

// SetAddSelection signals the type that we are adding to an
// existing selection instead of creating a new one.
func (c *CellSelector) SetAddSelection(v bool) {
	c.addSelection = v
}

// SelectionClear clears the current selection.
func (c *CellSelector) SelectionClear() {
	c.selection = c.selection[:0]
	c.selectionStart = nil
	c.selectionChanged = true
}

// SelectAll selects all cells.
func (c *CellSelector) SelectAll() {
	cells := sim.Cells()

	c.selectionStart = nil
	c.selection = make(sim.CellList, len(cells))
	c.selectionChanged = true

	copy(c.selection, cells)
	c.finalizeSelection()
}

// SelectionDelete deletes the current selection from the simulation.
func (c *CellSelector) SelectionDelete() {
	sel := c.selection
	if len(sel) == 0 {
		return
	}

	sim.Unload(c.selection)
	c.SelectionClear()
}

// SelectionMove moves the currently selected area up/down, left
// or right by some amount of cells as denoted by x and y. These
// can be be negative for left/up, positive for right/down and zero
// to not move at all.
func (c *CellSelector) SelectionMove(x, y int) {
	sel := c.selection
	if (x == 0 && y == 0) || len(sel) == 0 {
		return
	}

	// Remove the selection fom the simulation.
	sim.Unload(sel)

	// Move the selected nodes by the given offset.
	for i := 0; i < len(sel)-2; i += 3 {
		sel[i+0] += int32(x)
		sel[i+1] += int32(y)
	}

	// Paste the modified selection.
	sim.Load(0, 0, sel)
	sim.Trim()

	c.finalizeSelection()
}

func (c *CellSelector) MouseMove(x, y float64) {
	c.Grid.MouseMove(x, y)

	// If the viewport is being panned and we are defining a
	// selection area, move the starting point along with the
	// panned viewport. This allows us to select an area larger
	// than the viewport.
	if c.selecting && c.panning && c.selectionStart != nil {
		c.selectionStart.X -= int(c.mouseDelta[0])
		c.selectionStart.Y -= int(c.mouseDelta[1])
	}
}

func (c *CellSelector) MouseButton(button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	c.selecting = (button == glfw.MouseButton2)

	if action == glfw.Press {
		// Begin a new selection rectangle by storing the current
		// mouse position. This will form one of the corners of
		// the selection.
		if c.selecting {
			c.selectionStart = &image.Point{
				int(c.mousePosition[0]),
				int(c.mousePosition[1]),
			}
		}
		return
	}

	// No need to continue if we are not currently defining a selection.
	if !c.selecting || c.selectionStart == nil {
		return
	}

	// Find the selection rectangle.
	sr := c.selectionRect(*c.selectionStart, image.Point{
		int(c.mousePosition[0]),
		int(c.mousePosition[1]),
	})

	// Reset the selection area.
	c.selectionStart = nil

	// Find all cells in the selection rectangle.
	if !c.addSelection {
		c.selection = c.cellsInArea(sr)
		c.finalizeSelection()
		return
	}

	// Add the selected cells to the existing selection,
	// while making sure we don't have any duplicates.
	set := c.cellsInArea(sr)
	for i := 0; i < len(set)-2; i += 3 {
		x, y, v := set[i], set[i+1], set[i+2]

		if !c.selection.Contains(x, y) {
			c.selection = append(c.selection, x, y, v)
		}
	}

	c.finalizeSelection()
}

// finalizeSelection ensures there are no empty cells in the selection
// and it sorts the final selection.
func (c *CellSelector) finalizeSelection() {
	c.selection = c.selection.Trim()
	c.selectionChanged = true
}

// selectionRect computes and returns the canonical rectangle
// encompassing all selected grid cells.
func (c *CellSelector) selectionRect(ra, rb image.Point) image.Rectangle {
	or := image.Point{
		int(c.origin[0]),
		int(c.origin[1]),
	}

	cs := c.Zoom()
	rs := ra.Sub(or)
	re := rb.Sub(or)

	x1 := util.Min(rs.X, re.X) / cs
	y1 := util.Min(rs.Y, re.Y) / cs
	x2 := util.Max(rs.X, re.X) / cs
	y2 := util.Max(rs.Y, re.Y) / cs

	return image.Rect(x1, y1, x2, y2)
}

// cellsInArea finds all cells partially or entirely overlapping
// the given rectangle. It fills the out array with their
// respective indices.
func (c *CellSelector) cellsInArea(ra image.Rectangle) sim.CellList {
	cd := sim.Cells()
	out := make(sim.CellList, 0, 32)

	for i := 0; i < len(cd)-2; i += 3 {
		x := cd[i+0]
		y := cd[i+1]
		v := cd[i+2]

		if int(x) >= ra.Min.X && int(y) >= ra.Min.Y &&
			int(x) <= ra.Max.X && int(y) <= ra.Max.Y {
			out = append(out, x, y, v)
		}
	}

	return out
}

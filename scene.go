package main

import (
	"wireworld/resources"
	"wireworld/sim"
	"wireworld/ui"
	"wireworld/util"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Scene defines a window with loads of drawing and simulation
// manipulation functionality.
type Scene struct {
	projection  *util.Mat4
	window      *ui.Window
	canvas      *ui.Clipboard
	panel       *ui.InfoPanel
	history     sim.History
	currentTool int
	lmbPressed  bool
	infoVisible bool
}

// CreateScene creates a enw scene.
func CreateScene(c *Config) (*Scene, error) {
	var err error
	var s Scene

	s.window, err = ui.CreateWindow(int(c.Width), int(c.Height), c.Fullscreen)
	if err != nil {
		s.Release()
		return nil, err
	}

	err = resources.Load()
	if err != nil {
		s.Release()
		return nil, err
	}

	s.panel = ui.NewInfoPanel()
	s.canvas = ui.NewClipboard()
	s.currentTool = sim.CellWire
	s.lmbPressed = false
	s.infoVisible = true

	s.window.SetTitle(Version())
	s.window.SetKeyCallback(s.keyCallback)
	s.window.SetFramebufferSizeCallback(s.resizeCallback)
	s.window.SetMouseButtonCallback(s.mouseButtonCallback)
	s.window.SetCursorPosCallback(s.mouseMoveCallback)
	s.window.SetCharCallback(s.charCallback)
	s.window.SetScrollCallback(s.scrollCallback)

	// Make sure all components are initialized to the
	// correct dimensions.
	w, h := s.window.GetFramebufferSize()
	s.resizeCallback(nil, w, h)

	s.panel.Clear()
	return &s, nil
}

func (s *Scene) Release() {
	s.panel.Release()
	resources.Release()
	s.window.Release()
}

// Update updates the scene.
func (s *Scene) Update() bool {
	ok := s.window.Update()

	sim.Step(false)

	s.canvas.SetPanning(s.window.GetKey(glfw.KeySpace) == glfw.Press)
	s.updateInfo()
	return ok
}

// Draw renders the scene.
func (s *Scene) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	s.canvas.Draw(s.projection)

	if s.infoVisible {
		s.panel.Draw(s.projection)
	}
}

// drawCells draws on the grid. What is being drawn depends on the current mode.
func (s *Scene) drawCells() {
	if s.lmbPressed {
		x, y := s.canvas.HoverTarget()
		sim.Set(x, y, int32(s.currentTool))
	}
}

// setTool sets the current drawing tool.
func (s *Scene) setTool(t int) {
	s.currentTool = t
}

func (s *Scene) scrollCallback(_ *glfw.Window, x, y float64) {
	s.canvas.Scroll(x, y)
}

func (s *Scene) mouseMoveCallback(_ *glfw.Window, x, y float64) {
	s.canvas.MouseMove(x, y)
	s.drawCells()
}

func (s *Scene) mouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	s.canvas.MouseButton(button, action, mod)
	s.lmbPressed = (button == glfw.MouseButton1 && action == glfw.Press)
	s.drawCells()
}

func (s *Scene) resizeCallback(_ *glfw.Window, w, h int) {
	gl.Viewport(0, 0, int32(w), int32(h))

	s.projection = util.Mat4Ortho(0, float32(w), 0, float32(h), -1, 1)
	s.canvas.Resize(w, h)
	s.panel.Resize(0, 0, util.Max(w/5, 280), h)
}

func (s *Scene) charCallback(_ *glfw.Window, char rune) {

}

func (s *Scene) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Release:
		s.keyRelease(key, scancode, mods)
	case glfw.Press:
		s.keyPress(key, scancode, mods)
	}
}

func (s *Scene) keyRelease(key glfw.Key, scancode int, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyLeftShift, glfw.KeyRightShift:
		s.canvas.SetAddSelection(false)
	}
}

func (s *Scene) keyPress(key glfw.Key, scancode int, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyEscape:
		s.canvas.SelectionClear()
		s.canvas.ClipboardClear()
	case glfw.KeyF1:
		s.canvas.ToggleGridVisible()
	case glfw.KeyF2:
		s.canvas.ToggleDrawClipboard()

	case glfw.KeyGraveAccent:
		s.infoVisible = !s.infoVisible

	case glfw.KeyQ:
		sim.ToggleRunning()
	case glfw.KeyE:
		sim.Step(true)
	case glfw.KeyT:
		sim.Trim()

	case glfw.Key1:
		s.setTool(sim.CellEmpty)
	case glfw.Key2:
		s.setTool(sim.CellWire)
	case glfw.Key3:
		s.setTool(sim.CellHead)
	case glfw.Key4:
		s.setTool(sim.CellTail)

	case glfw.KeyEqual:
		sim.ScaleInterval(-1)
	case glfw.KeyMinus:
		sim.ScaleInterval(+1)

	case glfw.KeyLeftShift, glfw.KeyRightShift:
		s.canvas.SetAddSelection(true)
	case glfw.KeyA:
		if mods&glfw.ModControl != 0 {
			s.canvas.SelectAll()
		}
	case glfw.KeyDelete:
		s.canvas.SelectionDelete()
	case glfw.KeyUp:
		s.canvas.SelectionMove(0, -1)
	case glfw.KeyDown:
		s.canvas.SelectionMove(0, +1)
	case glfw.KeyLeft:
		s.canvas.SelectionMove(-1, 0)
	case glfw.KeyRight:
		s.canvas.SelectionMove(+1, 0)

	case glfw.KeyC:
		if mods&glfw.ModControl != 0 {
			s.canvas.ClipboardCopy()
		}
	case glfw.KeyV:
		if mods&glfw.ModControl != 0 {
			s.canvas.ClipboardPaste()
		} else {
			s.canvas.ScrollTo(0, 0)
			s.canvas.SetZoom(ui.ZoomDefault)
		}
	case glfw.KeyX:
		if mods&glfw.ModControl != 0 {
			s.canvas.ClipboardCut()
		}

	case glfw.KeyZ:
		if mods&glfw.ModControl != 0 && mods&glfw.ModShift == 0 {
			s.history.Undo()
		} else if mods&glfw.ModControl != 0 && mods&glfw.ModShift != 0 {
			s.history.Redo()
		}
	}
}

// updateInfo recreates the text contents of the debug/info panel.
func (s *Scene) updateInfo() {
	if !s.infoVisible {
		return
	}

	line := 1
	s.panel.Clear()
	p := func(v string, argv ...interface{}) {
		s.panel.Print(line, v, argv...)
		line++
	}

	p("Cells: %d, running: %v", sim.CellCount(), sim.Running())
	p("Step interval: %s", sim.StepInterval())
	p("Current tool: %s", toolName(s.currentTool))

	p("")
	p("Simulation:")
	p(" [q] Start/stop simulation")
	p(" [e] Single simulation step")
	p(" [+] Double simulation speed")
	p(" [-] Halve simulation speed")

	p("")
	p("Tools:")
	p(" [1] Draw Empty cell")
	p(" [2] Draw Wire cell")
	p(" [3] Draw Electron head")
	p(" [4] Draw Electron tail")
	p(" [t] Trim empty cells")
	p(" [ctrl-a] Select all cells")
	p(" [ctrl-x] Cut selection")
	p(" [ctrl-c] Copy selection")
	p(" [ctrl-v] Paste selection")
	p(" [del] Delete selection")
	p(" [arrow keys] Move selection")

	p("")
	p("Misc:")
	p(" [~] Show/hide this info panel")
	p(" [F1] Toggle grid visibility")
	p(" [F2] Toggle clipboard visibility")
	p(" [esc] Cancel selection / Clear clipboard")
	p(" [lmb] Draw cells")
	p(" [rmb] Draw selection")
	p(" [wheel] Zoom in/out")
	p(" [space+mouse] Pan viewport")
}

// toolName returns a human-readable name for the given cell state.
func toolName(v int) string {
	switch v {
	case sim.CellWire:
		return "Wire"
	case sim.CellHead:
		return "Electron head"
	case sim.CellTail:
		return "Electron tail"
	default:
		return "Empty"
	}
}

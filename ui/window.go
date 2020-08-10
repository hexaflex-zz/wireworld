package ui

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Canvas creates a GLFW and OpenGL context and facilitates
// handling of user input.
type Window struct {
	*glfw.Window
}

// CreateWindow creates a new window with the given configuration.
func CreateWindow(width, height int, fullscreen bool) (*Window, error) {
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	var w Window

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if fullscreen {
		mon := glfw.GetPrimaryMonitor()
		mode := mon.GetVideoMode()
		w.Window, err = glfw.CreateWindow(mode.Width, mode.Height, "", mon, nil)
	} else {
		w.Window, err = glfw.CreateWindow(width, height, "", nil, nil)
	}

	if err != nil {
		glfw.Terminate()
		return nil, err
	}

	w.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		w.Destroy()
		glfw.Terminate()
		return nil, err
	}

	glfw.SwapInterval(1)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	return &w, nil
}

// Release releases all GLFW and OpenGL resources.
func (w *Window) Release() {
	if w.Window != nil {
		w.SetKeyCallback(nil)
		w.SetFramebufferSizeCallback(nil)
		w.SetFocusCallback(nil)
		w.SetMouseButtonCallback(nil)
		w.SetCursorPosCallback(nil)
		w.SetCharCallback(nil)
		w.SetScrollCallback(nil)
		w.SetUserPointer(nil)
		w.Destroy()
		w.Window = nil
	}

	glfw.Terminate()
}

// SetTitle sets the window title.
func (w *Window) SetTitle(f string, argv ...interface{}) {
	w.Window.SetTitle(fmt.Sprintf(f, argv...))
}

// Update calls PollEvents, swaps the buffers and returns true
// if the window is to be closed.
func (w *Window) Update() bool {
	w.Window.SwapBuffers()
	glfw.PollEvents()
	return !w.ShouldClose()
}

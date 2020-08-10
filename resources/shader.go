package resources

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Shader loads GLSL shader sources and facilitates setting
// uniform values.
type Shader struct {
	name string // Name, used in error messages.
	id   uint32 // Program id.
}

// CompileShader loads the given vertex-, geometry- and fragment shader
// sources and compiles them into a single program. The name value is
// used in error emssages to more easily identify the source of a problem.
func CompileShader(name, vsSrc, gsSrc, fsSrc string) (*Shader, error) {
	vsSrc = strings.TrimSpace(vsSrc)
	gsSrc = strings.TrimSpace(gsSrc)
	fsSrc = strings.TrimSpace(fsSrc)

	if len(vsSrc) == 0 && len(gsSrc) == 0 && len(fsSrc) == 0 {
		return nil, errors.New(name + ": no shader sources defined")
	}

	vs, err := compile(name, vsSrc, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	gs, err := compile(name, gsSrc, gl.GEOMETRY_SHADER)
	if err != nil {
		gl.DeleteShader(vs)
		return nil, err
	}

	fs, err := compile(name, fsSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		gl.DeleteShader(vs)
		gl.DeleteShader(gs)
		return nil, err
	}

	return link(name, vs, fs, gs)
}

func (s *Shader) Release() {
	gl.UseProgram(0)
	gl.DeleteProgram(s.id)
}

func (s *Shader) Use()   { gl.UseProgram(s.id) }
func (s *Shader) Unuse() { gl.UseProgram(0) }

// set finds a named uniform location and calls f with its location value.
// Returns an error if the location can not be found.
func (s *Shader) set(name string, f func(int32)) error {
	loc := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	if loc < 0 {
		err := fmt.Errorf("%s: uniform %q not found", s.name, name)
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	f(loc)
	return nil
}

// Set1i sets the given uniform.
func (s *Shader) Set1i(name string, v int) error {
	return s.set(name, func(loc int32) {
		gl.Uniform1i(loc, int32(v))
	})
}

// Set1f sets the given uniform.
func (s *Shader) Set1f(name string, v float32) error {
	return s.set(name, func(loc int32) {
		gl.Uniform1f(loc, v)
	})
}

// Set2f sets the given uniform.
func (s *Shader) Set2f(name string, a, b float32) error {
	return s.set(name, func(loc int32) {
		gl.Uniform2f(loc, a, b)
	})
}

// Set4f sets the given uniform.
func (s *Shader) Set4f(name string, a, b, c, d float32) error {
	return s.set(name, func(loc int32) {
		gl.Uniform4f(loc, a, b, c, d)
	})
}

// SetMat16 sets the given uniform.
func (s *Shader) SetMat16(name string, m []float32) error {
	return s.set(name, func(loc int32) {
		gl.UniformMatrix4fv(loc, 1, false, &m[0])
	})
}

// compile compiles a specific shader.
func compile(name, source string, stype uint32) (uint32, error) {
	if len(source) == 0 {
		return 0, nil
	}

	s := gl.CreateShader(stype)
	cstr, free := gl.Strs(source + "\x00")
	gl.ShaderSource(s, 1, cstr, nil)
	free()
	gl.CompileShader(s)

	// Make sure compilation worked.
	var value int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &value)

	if value == gl.TRUE {
		return s, nil
	}

	// Compilation failed. We want to know why.
	gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &value)

	log := gl.Str(strings.Repeat("\x00", int(value+1)))
	gl.GetShaderInfoLog(s, value, nil, log)
	gl.DeleteShader(s)

	var sname string
	switch stype {
	case gl.VERTEX_SHADER:
		sname = "vertex"
	case gl.GEOMETRY_SHADER:
		sname = "geometry"
	case gl.FRAGMENT_SHADER:
		sname = "fragment"
	}

	return 0, fmt.Errorf("%s (%s): %s", name, sname, gl.GoStr(log))
}

// link links a shader program. The name is used in error messages if
// applicable.
func link(name string, vs, gs, fs uint32) (*Shader, error) {
	p := gl.CreateProgram()

	if vs > 0 {
		gl.AttachShader(p, vs)
	}

	if gs > 0 {
		gl.AttachShader(p, gs)
	}

	if fs > 0 {
		gl.AttachShader(p, fs)
	}

	gl.LinkProgram(p)

	if vs > 0 {
		gl.DetachShader(p, vs)
		gl.DeleteShader(vs)
	}

	if gs > 0 {
		gl.DetachShader(p, gs)
		gl.DeleteShader(gs)
	}

	if fs > 0 {
		gl.DetachShader(p, fs)
		gl.DeleteShader(fs)
	}

	var value int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &value)
	if value == gl.TRUE {
		return &Shader{name, p}, nil
	}

	// Linking failed. Find out why and clean up the mess.
	gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &value)

	log := gl.Str(strings.Repeat("\x00", int(value+1)))
	gl.GetProgramInfoLog(p, value, nil, log)
	gl.DeleteProgram(p)

	return nil, fmt.Errorf("%s (linking): %s", name, gl.GoStr(log))
}

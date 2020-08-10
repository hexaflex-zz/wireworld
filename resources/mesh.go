package resources

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Mesh defines a generic mesh object.
type Mesh interface {
	Release()
	Draw()

	Commitfv([]float32, uint32)
	Commitiv([]int32, uint32)
}

// TexturedQuadMesh defines a simple quad.
type TexturedQuadMesh struct {
	vao     uint32
	vboPos  uint32
	vboTex  uint32
	texture uint32
}

func newTexturedQuadMesh() Mesh {
	var m TexturedQuadMesh

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	quad := []float32{0, 0, 1, 0, 0, 1, 1, 1}

	gl.GenBuffers(1, &m.vboPos)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vboPos)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.BufferData(gl.ARRAY_BUFFER, 8*4, gl.Ptr(quad), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &m.vboTex)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vboTex)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)
	gl.BufferData(gl.ARRAY_BUFFER, 8*4, gl.Ptr(quad), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)

	gl.GenTextures(1, &m.texture)
	return &m
}

// Release clears mesh resources.
func (m *TexturedQuadMesh) Release() {
	gl.DeleteTextures(1, &m.texture)
	gl.DeleteBuffers(1, &m.vboPos)
	gl.DeleteBuffers(1, &m.vboTex)
	gl.DeleteVertexArrays(1, &m.vao)
}

// Draw renders the mesh.
func (m *TexturedQuadMesh) Draw() {
	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, m.texture)

	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (m *TexturedQuadMesh) Commitiv([]int32, uint32)   {}
func (m *TexturedQuadMesh) Commitfv([]float32, uint32) {}

// CommitTexture uploads the given texture and sets UV coordinates
// for the given size. This size is not necessarily the same as
// the size of img. Meaning w/h define a subset of img.
func (m *TexturedQuadMesh) CommitTexture(img *image.RGBA, w, h int) {
	tb := img.Bounds()

	// Upload the texture data.
	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, m.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(tb.Dx()), int32(tb.Dy()),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	// Define the UV coordinates.
	//
	// We divide all values by the texture dimensions to convert
	// the pixel offsets to UV coordinates in the range [0,1].
	rx1 := float32(0) // x / float32(tb.Dx())
	ry1 := float32(0) // y / float32(tb.Dy())
	rx2 := float32(w) / float32(tb.Dx())
	ry2 := float32(h) / float32(tb.Dy())

	uv := []float32{rx1, ry1, rx2, ry1, rx1, ry2, rx2, ry2}
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vboTex)
	gl.BufferData(gl.ARRAY_BUFFER, 8*4, gl.Ptr(uv), gl.STREAM_DRAW)
}

// QuadMesh defines a simple, untextured quad.
type QuadMesh struct {
	vao uint32
	vbo uint32
}

func newQuadMesh() Mesh {
	var m QuadMesh

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	quad := []float32{0, 0, 1, 0, 0, 1, 1, 1}

	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.BufferData(gl.ARRAY_BUFFER, 8*4, gl.Ptr(&quad[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)

	return &m
}

// Release clears mesh resources.
func (m *QuadMesh) Release() {
	gl.DeleteVertexArrays(1, &m.vao)
	gl.DeleteBuffers(1, &m.vbo)
}

// Draw renders the mesh.
func (m *QuadMesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (m *QuadMesh) Commitiv([]int32, uint32)   {}
func (m *QuadMesh) Commitfv([]float32, uint32) {}

// GridMesh defines a grid.
type GridMesh struct {
	vao  uint32
	vbo  uint32
	size int32
}

func newGridMesh() Mesh {
	var m GridMesh
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	return &m
}

// Release clears mesh resources.
func (m *GridMesh) Release() {
	gl.DeleteVertexArrays(1, &m.vao)
	gl.DeleteBuffers(1, &m.vbo)
}

// Draw renders the mesh.
func (m *GridMesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.LINES, 0, m.size)
}

func (m *GridMesh) Commitiv([]int32, uint32) {}
func (m *GridMesh) Commitfv(set []float32, usage uint32) {
	m.size = int32(len(set) / 2)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(set)*4, gl.Ptr(set), usage)
}

// CellMesh defines a set of simulation cells.
type CellMesh struct {
	vao  uint32
	vbo  uint32
	size int32
}

func newCellMesh() Mesh {
	var m CellMesh
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.VertexAttribIPointer(0, 3, gl.INT, 3*4, nil)
	gl.EnableVertexAttribArray(0)
	return &m
}

// Release clears mesh resources.
func (m *CellMesh) Release() {
	gl.DeleteVertexArrays(1, &m.vao)
	gl.DeleteBuffers(1, &m.vbo)
}

// Draw renders the mesh.
func (m *CellMesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.POINTS, 0, m.size)
}

func (m *CellMesh) Commitfv([]float32, uint32) {}
func (m *CellMesh) Commitiv(set []int32, usage uint32) {
	m.size = int32(len(set) / 3)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	if len(set) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(set)*4, gl.Ptr(set), usage)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, 1, nil, usage)
	}
}

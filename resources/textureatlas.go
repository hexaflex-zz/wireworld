// https://github.com/go-gl-legacy/glh/blob/master/atlas.go

/*
Copyright (c) 2012 The go-gl Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
     notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
     copyright notice, this list of conditions and the following disclaimer
     in the documentation and/or other materials provided with the
     distribution.
   * Neither the name of go-gl nor the names of its contributors may be used
     to endorse or promote products derived from this software without specific
     prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package resources

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"sync"

	"wireworld/util"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// A node represents an area of an atlas texture which
// has been allocated for use.
type atlasNode struct {
	x int // region x
	y int // region y + height
	z int // region width
}

// A texture atlas is used to tightly pack arbitrarily many small images
// into a single texture.
//
// The actual implementation is based on the article by Jukka Jylänki:
// "A Thousand Ways to Pack the Bin - A Practical Approach to Two-Dimensional
// Rectangle Bin Packing", February 27, 2010.
//
// More precisely, this is an implementation of the
// 'Skyline Bottom-Left' algorithm.
type TextureAtlas struct {
	m              sync.RWMutex
	nodes          []atlasNode // Allocated nodes.
	data           *image.RGBA // Atlas pixel data.
	texture        uint32      // Glyph texture.
	textureChanged bool        // texture has changed?
}

// NewAtlas creates a new texture atlas.
//
// The given width, height determine the size of the underlying texture.
// The values are scaled up to the nearest power-of-two value, if they are
// not already a power-of-two.
func NewTextureAtlas(width, height int) *TextureAtlas {
	ta := new(TextureAtlas)
	ta.textureChanged = true
	ta.data = image.NewRGBA(image.Rect(0, 0, util.Pow2(width), util.Pow2(height)))

	gl.GenTextures(1, &ta.texture)

	// We want a one pixel border around the whole atlas to avoid
	// any artefacts when sampling our texture.
	ta.nodes = []atlasNode{{1, 1, width - 2}}

	return ta
}

// Release clears all atlas resources.
func (ta *TextureAtlas) Release() {
	if ta == nil {
		return
	}

	ta.m.Lock()
	gl.DeleteTextures(1, &ta.texture)
	ta.data = nil
	ta.nodes = nil
	ta.m.Unlock()
}

// Size returns the texture dimensions.
func (ta *TextureAtlas) Size() (int, int) {
	ta.m.Lock()
	b := ta.data.Bounds()
	ta.m.Unlock()
	return b.Dx(), b.Dy()
}

// Clear removes all allocated regions from the atlas.
// This invalidates any previously allocated regions.
func (ta *TextureAtlas) Clear() {
	ta.m.Lock()
	ta.textureChanged = true
	ta.nodes = ta.nodes[:1]

	tw := ta.data.Bounds().Dx()

	// We want a one pixel border around the whole atlas to avoid
	// any artefacts when sampling our texture.
	ta.nodes[0].x = 1
	ta.nodes[0].y = 1
	ta.nodes[0].z = tw - 2

	draw.Draw(ta.data, ta.data.Bounds(),
		image.Transparent, image.ZP, draw.Src)

	ta.m.Unlock()
}

// Bind binds the atlas texture, so it can be used for rendering.
// This implictely commits the texture data to the GPU if it has
// changed since the last Bind call.
func (ta *TextureAtlas) Bind() {
	ta.m.RLock()

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, ta.texture)

	if ta.textureChanged {
		ta.commit()
		ta.textureChanged = false
	}

	ta.m.RUnlock()
}

// Unbind unbinds the current texture.
// Note that this applies to any texture currently active.
// If this is not the atlas texture, it will still perform the action.
func (ta *TextureAtlas) Unbind() {
	ta.m.RLock()
	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	ta.m.RUnlock()
}

// commit creates the actual texture from the atlas image data.
// This should be called after all regions have been defined and set,
// and before you start using the texture for display.
func (ta *TextureAtlas) commit() {
	b := ta.data.Bounds()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(b.Dx()), int32(b.Dy()),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(ta.data.Pix))
}

// Allocate allocates a new region of the given dimensions in the atlas.
// It returns false if the allocation failed. This can happen when the
// specified dimensions exceed atlas bounds, or the atlas is full.
func (ta *TextureAtlas) Allocate(width, height int) (image.Rectangle, bool) {
	ta.m.Lock()
	defer ta.m.Unlock()

	region := image.Rect(0, 0, width, height)
	bestIndex := -1
	bestWidth := 1<<31 - 1
	bestHeight := 1<<31 - 1

	for index := range ta.nodes {
		y := ta.fit(index, width, height)

		if y < 0 {
			continue
		}

		node := ta.nodes[index]

		if ((y + height) < bestHeight) || (((y + height) == bestHeight) && (node.z < bestWidth)) {
			bestWidth = node.z
			bestHeight = y + height
			bestIndex = index
			region = image.Rect(node.x, y, node.x+width, y+height)
		}
	}

	// There is no suitable space for the new allocation.
	if bestIndex == -1 {
		return region, false
	}

	// Insert a new node at bestIndex
	ta.nodes = append(ta.nodes, atlasNode{})
	copy(ta.nodes[bestIndex+1:], ta.nodes[bestIndex:])
	ta.nodes[bestIndex] = atlasNode{region.Min.X, region.Min.Y + height, width}

	// Adjust subsequent nodes.
	for i := bestIndex + 1; i < len(ta.nodes); i++ {
		curr := &ta.nodes[i]
		prev := &ta.nodes[i-1]

		if curr.x >= prev.x+prev.z {
			break
		}

		shrink := prev.x + prev.z - curr.x
		curr.x += shrink
		curr.z -= shrink

		if curr.z > 0 {
			break
		}

		copy(ta.nodes[i:], ta.nodes[i+1:])
		ta.nodes = ta.nodes[:len(ta.nodes)-1]
		i--
	}

	ta.merge()
	return region.Canon(), true
}

// Set draws the given image into the given region of the atlas.
// It assumes there is enough space available for the data to fit.
func (ta *TextureAtlas) Set(region image.Rectangle, src *image.RGBA) {
	ta.m.RLock()
	draw.Draw(ta.data, region, src, src.Bounds().Min, draw.Src)
	ta.textureChanged = true
	ta.m.RUnlock()
}

// fit checks if the given dimensions fit in the given node.
// If not, it checks any subsequent nodes for a match.
// It returns the height for the last checked node.
// Returns -1 if the width or height exceed texture capacity.
func (ta *TextureAtlas) fit(index, width, height int) int {
	node := ta.nodes[index]
	b := ta.data.Bounds()
	tw, th := b.Dx(), b.Dy()

	if node.x+width > tw-1 {
		return -1
	}

	y := node.y
	remainder := width

	for remainder > 0 {
		node = ta.nodes[index]

		if node.y > y {
			y = node.y
		}

		if y+height > th-1 {
			return -1
		}

		remainder -= node.z
		index++
	}

	return y
}

// merge merges nodes where possible.
// This is the case when two regions overlap.
func (ta *TextureAtlas) merge() {
	for i := 0; i < len(ta.nodes)-1; i++ {
		node := &ta.nodes[i]
		next := ta.nodes[i+1]

		if node.y != next.y {
			continue
		}

		node.z += next.z

		copy(ta.nodes[i+1:], ta.nodes[i+2:])
		ta.nodes = ta.nodes[:len(ta.nodes)-1]
		i--
	}
}

// Save saves the texture image as a PNG file with the given name.
func (a *TextureAtlas) Save(file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	return png.Encode(fd, a.data)
}

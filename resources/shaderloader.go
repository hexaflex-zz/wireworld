package resources

import (
	"strings"
	"sync"
)

// shaders provides global access to shader resources.
var shaders = newShaderLoader()

// GetShader returns the shader for the given name.
func GetShader(name string) *Shader {
	return shaders.Get(name)
}

// shaderLoader loads all necessary shader programs.
type shaderLoader struct {
	m        sync.RWMutex
	programs map[string]*Shader
}

// newShaderLoader creates a new, empty shader manager.
func newShaderLoader() *shaderLoader {
	return &shaderLoader{
		programs: make(map[string]*Shader),
	}
}

// Release clears all shader resources.
func (sl *shaderLoader) Release() {
	sl.m.Lock()

	for k, v := range sl.programs {
		v.Release()
		delete(sl.programs, k)
	}

	sl.m.Unlock()
}

// Load loads all known shader programs.
func (sl *shaderLoader) Load() error {
	sl.m.Lock()
	defer sl.m.Unlock()

	if err := sl.loadShader("CellRenderer", cellRendererSources); err != nil {
		return err
	}

	if err := sl.loadShader("CellSelectorRect", cellSelectorRectSources); err != nil {
		return err
	}

	if err := sl.loadShader("CellSelectorCells", cellSelectorCellsSources); err != nil {
		return err
	}

	if err := sl.loadShader("Grid", gridSources); err != nil {
		return err
	}

	if err := sl.loadShader("Panel", panelSources); err != nil {
		return err
	}

	return nil
}

// Get returns the shader associated with the given name.
func (sl *shaderLoader) Get(name string) *Shader {
	sl.m.RLock()
	v := sl.programs[strings.ToLower(name)]
	sl.m.RUnlock()
	return v
}

// loadShader loads a specific shader program.
func (sl *shaderLoader) loadShader(name string, src [3]string) error {
	var err error
	name = strings.ToLower(name)
	sl.programs[name], err = CompileShader(name, src[0], src[1], src[2])
	return err
}

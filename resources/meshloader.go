package resources

import (
	"strings"
	"sync"
)

// meshes provides global access to mesh resources.
var meshes = newMeshLoader()

// GetMesh returns the mesh for the given name.
func GetMesh(name string) Mesh {
	return meshes.Get(name)
}

// meshLoader loads all necessary mesh objects.
type meshLoader struct {
	m   sync.RWMutex
	set map[string]Mesh
}

// newMeshLoader creates a new, empty mesh manager.
func newMeshLoader() *meshLoader {
	return &meshLoader{
		set: make(map[string]Mesh),
	}
}

// Release clears all shader resources.
func (ml *meshLoader) Release() {
	ml.m.Lock()

	for k, v := range ml.set {
		v.Release()
		delete(ml.set, k)
	}

	ml.m.Unlock()
}

// Load loads all known mesh objects.
func (ml *meshLoader) Load() error {
	ml.m.Lock()

	ml.loadMesh("Panel", newTexturedQuadMesh())
	ml.loadMesh("CellSelectorRect", newQuadMesh())
	ml.loadMesh("CellSelectorCells", newCellMesh())
	ml.loadMesh("CellRenderer", newCellMesh())
	ml.loadMesh("Clipboard", newCellMesh())
	ml.loadMesh("Grid", newGridMesh())

	ml.m.Unlock()
	return nil
}

// Get returns the mesh associated with the given name.
func (ml *meshLoader) Get(name string) Mesh {
	ml.m.RLock()
	v := ml.set[strings.ToLower(name)]
	ml.m.RUnlock()
	return v
}

// loadMesh loads a specific mesh.
func (ml *meshLoader) loadMesh(name string, m Mesh) {
	ml.set[strings.ToLower(name)] = m
}

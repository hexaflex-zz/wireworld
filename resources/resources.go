// Package resources defines and loads shader and mesh resources.
package resources

// Load loads all resources.
func Load() error {
	err := shaders.Load()
	if err != nil {
		return err
	}

	err = meshes.Load()
	if err != nil {
		Release()
		return err
	}

	return nil
}

// Release clears all resources.
func Release() {
	shaders.Release()
	meshes.Release()
}

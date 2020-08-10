package main

import (
	"fmt"
	"os"
	"runtime"
)

// Make sure main() and all openGL related stuff runs in
// the main thread.r
func init() { runtime.LockOSThread() }

func main() {
	config := ParseArgs()

	// Initialize the window, opengl and all scene related things.
	scene, err := CreateScene(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for scene.Update() {
		scene.Draw()
	}

	scene.Release()
}

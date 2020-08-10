package main

import (
	"flag"
	"fmt"
	"os"
)

// Config defines application settings.
type Config struct {
	Width      uint
	Height     uint
	Fullscreen bool
}

// ParseArgs parses commandline arguments and returns a config struct.
// Exits the program with an error if invalid data was found.
func ParseArgs() *Config {
	var c Config
	c.Width = 1280
	c.Height = 800
	c.Fullscreen = false

	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.UintVar(&c.Width, "width", c.Width, "Display width in pixels.")
	flag.UintVar(&c.Height, "height", c.Height, "Display height in pixels.")
	flag.BoolVar(&c.Fullscreen, "fullscreen", c.Fullscreen, "Use a fullscreen or windowed display.")
	version := flag.Bool("version", false, "Displays version information.")
	flag.Parse()

	if *version {
		fmt.Println(Version())
		os.Exit(0)
	}

	if c.Width == 0 {
		fmt.Fprintf(os.Stderr, "width should be > 0")
		flag.Usage()
		os.Exit(1)
	}

	if c.Height == 0 {
		fmt.Fprintf(os.Stderr, "height should be > 0")
		flag.Usage()
		os.Exit(1)
	}

	return &c
}

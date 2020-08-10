package main

import (
	"fmt"
)

const (
	AppName         = "wireworld"
	AppVersionMajor = 0
	AppVersionMinor = 21
)

func Version() string {
	return fmt.Sprintf("%s v%d.%d\n",
		AppName, AppVersionMajor, AppVersionMinor)
}

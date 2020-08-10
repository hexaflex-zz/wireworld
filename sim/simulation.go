// Package sim implements the data and rules for the Wireworld
// cellular automata.
//
// ref: https://en.wikipedia.org/wiki/Wireworld
package sim

import "time"

// Known cell states.
const (
	CellEmpty = iota
	CellWire
	CellHead
	CellTail
)

var (
	data simulationData

	// stepInterval defines the time, in microseconds, between each
	// step cycle if the simulation is running.
	stepInterval = 50 * time.Millisecond
	stepTimer    time.Time

	// running determines if the simulation is currently running by itself.
	running bool
)

// CellsChanged returns true if the cell buffer has changed since
// the last call to CellsChanged. This call implicitely resets the
// cellsChanged flag.
func CellsChanged() bool {
	ok := data.cellsChanged
	data.cellsChanged = false
	return ok
}

// CellCount returns the number of cells in the simulation.
func CellCount() int {
	return data.CellCount()
}

// Cells returns the cell buffer.
func Cells() CellList {
	return data.cellData
}

// SetList replaces all cells present in the simulation and v with the value from v.
func SetList(v CellList) {
	data.UpdateList(v)
}

// StepInterval returns the current step interval.
func StepInterval() time.Duration {
	return stepInterval
}

// ScaleInterval sets the new step interval by halving or doubling the
// current value. There is a lower bound of 1 microsecond.
// Delta is expected to be -1 or +1.
func ScaleInterval(delta int) {
	v := stepInterval

	if delta < 0 {
		v = v >> 1
	} else {
		v = v << 1
	}

	// We don't want to go below 1 microsecond.
	if v < time.Microsecond {
		v = time.Microsecond
	}

	// Round the number up or down to something reasonable.
	switch {
	case v >= time.Second:
		v = v.Truncate(time.Second)
	case v >= time.Millisecond:
		v = v.Truncate(time.Millisecond)
	default:
		v = v.Truncate(time.Microsecond)
	}

	stepInterval = v
}

// Running returns the running state.
func Running() bool {
	return running
}

// ToggleRunning toggles the running state and returns the new state.
func ToggleRunning() bool {
	running = !running
	stepTimer = time.Now()
	return running
}

// Trim removes any cells with the CellEmpty value.
// These have no effect on the simulation and only take up space.
//
// Its use is mostly to squeeze every bit of speed out of the
// Step() function by eliminating unnecessary cells. Additionally
// to minimize the amount of memory used, so the simulation can
// be stored on disk efficiently.
func Trim() {
	data.Trim()
}

// Load loads the given set of cell states into the simulation.
// Its top-right corner is placed at the given position.
func Load(x, y int32, set CellList) {
	data.Load(x, y, set)
	data.Sort()
}

// Unload removes all the given cells from the simulation.
// This marks all existing cells as CellEmpty. To really
// delete them, use the Trim() function afterwards.
func Unload(set CellList) {
	data.Unload(set)
}

// Set sets the cell at position x/y to the given state.
//
// If the target cell does not yet exist and the state is CellEmpty,
// this call does nothing. If an existing cell is set to CellEmpty,
// it will not be deleted. The Trim() function is meant to do that
// whenever called separately.
func Set(x, y, state int32) {
	data.Set(x, y, state)
	data.Sort()
}

// Step applies the wireworld rules to the celldata once.
// If force is true, this is done immediately and unconditionally.
// If force is false, this call is ignored if not enough time has
// passed since the last step() call. 'Enough time' is determined by
// the value of stepInterval.
func Step(force bool) {
	// Make sure we are actually meant to perform the step call.
	if !force {
		if !running {
			return
		}

		now := time.Now()
		if now.Sub(stepTimer) < stepInterval {
			return
		}

		stepTimer = now
	}

	data.Step()
}

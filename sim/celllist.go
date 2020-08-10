package sim

import (
	"sort"
)

// CellList defines a set of cells.
type CellList []int32

// Trim removes any cells with the CellEmpty value and returns
// the resulting, new list. Empty cells have no effect on the
// simulation and only take up space.
//
// Its use is mostly to squeeze every bit of speed out of the
// Step() function by eliminating unnecessary cells. Additionally
// to minimize the amount of memory used, so the simulation can
// be stored on disk efficiently.
func (c CellList) Trim() CellList {
	newList := make(CellList, 0, len(c))

	for i := 0; i < len(c)-2; i += 3 {
		if c[i+2] != CellEmpty {
			newList = append(newList, c[i], c[i+1], c[i+2])
		}
	}

	return newList
}

// Load loads c2 to c1 and returns the resulting set.
// c2's top-right corner is placed at the given position.
// This ensures no cell duplicates are added.
func (c1 CellList) Load(x, y int32, c2 CellList) CellList {
	newList := c1

	// Add cells if they do not yet exist.
	for i := 0; i < len(c2)-2; i += 3 {
		// Move the new cell to its target location by
		// offsetting it by the provided x/y coordinates.
		cx := c2[i+0] + x
		cy := c2[i+1] + y
		cv := c2[i+2]

		// Find out of the cell already exists or not. We only search
		// through the existing cell data. That part is sorted, so
		// the indexOf call should work as intended.
		n := c1.IndexOf(cx, cy)
		if n > -1 {
			// Cell already exists. Update its value.
			newList[n+2] = cv
		} else {
			// Append new cell.
			newList = append(newList, cx, cy, cv)
		}
	}

	return newList
}

// Unload returns c1 where all entries from c2 have been removed from c1.
func (c1 CellList) Unload(c2 CellList) CellList {
	for i := 0; i < len(c2)-2; i += 3 {
		n := c1.IndexOf(c2[i+0], c2[i+1])
		if n > -1 {
			c1[n+2] = CellEmpty
		}
	}

	return c1.Trim()
}

// Len returns the number of cells in the list.
// This is not the same as the length of the CellList slice.
func (c CellList) Len() int { return len(c) / 3 }
func (c CellList) Sort()    { sort.Sort(c) }

func (c CellList) Swap(i, j int) {
	a, b := i*3, j*3
	c[a+0], c[b+0] = c[b+0], c[a+0]
	c[a+1], c[b+1] = c[b+1], c[a+1]
	c[a+2], c[b+2] = c[b+2], c[a+2]
}

func (c CellList) Less(i, j int) bool {
	a, b := i*3, j*3
	xa, xb := c[a+0], c[b+0]
	ya, yb := c[a+1], c[b+1]
	return (xa < xb) || ((xa == xb) && (ya < yb))
}

// Contains returns true if the cell with coordinates X/Y exists in the set.
func (c CellList) Contains(x, y int32) bool {
	return c.IndexOf(x, y) > -1
}

// IndexOf returns the index of the cell with cooridnates X/Y in c.
// Returns -1 if it can't be found.
//
// This performs a binary search and assumes c is sorted.
func (c CellList) IndexOf(x, y int32) int {
	if len(c) == 0 {
		return -1
	}

	var i, m, n int
	j := (len(c) / 3) - 1

	for i < j {
		m = int(uint(i+j) >> 1)
		n = m * 3

		if (c[n] < x) || (c[n] == x && c[n+1] < y) {
			i = m + 1
		} else {
			j = m
		}
	}

	n = i * 3
	if i != j || c[n] != x || c[n+1] != y {
		return -1
	}

	return n
}

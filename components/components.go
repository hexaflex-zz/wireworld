// Package components contains some predefined wireworld components
// which can be loaded into a simulation.
//
// Each component defines a list of cells. Each cell is composed
// of 3 integers: X and Y cell cooridnates and the cell state.
// Cell states can have one of the following values:
//
//   0: empty
//   1: wire
//   2: electron head
//   3: electron tail
//
// Empty cells (state 0) can be omitted altogether. They just take
// up unnecessary space.
//
// ref: https://www.quinapalus.com/wi-index.html
package components

// Clock4 defines a clock with a 4-cycle interval.
var Clock4 = []int32{1, 0, 2, 0, 1, 1, 2, 1, 3, 1, 2, 1}

// Diode defines a 1-way wire.
var Diode = []int32{
	3, 0, 1, 4, 0, 1,
	0, 1, 3, 1, 1, 2,
	2, 1, 1, 4, 1, 1,
	5, 1, 1, 6, 1, 1,
	3, 2, 1, 4, 2, 1,
}

// OR defines an OR gate.
var OR = []int32{
	0, 0, 1, 1, 0, 1, 2, 0, 1, 3, 1, 1,
	2, 2, 1, 3, 2, 1, 4, 2, 1, 5, 2, 1, 6, 2, 1,
	3, 3, 1, 0, 4, 1, 1, 4, 1, 2, 4, 1,
}

// XOR defines an exclusive-OR gate.
var XOR = []int32{
	0, 0, 1, 1, 0, 1, 2, 0, 1, 3, 1, 1,
	2, 2, 1, 3, 2, 1, 4, 2, 1, 5, 2, 1,
	2, 3, 1, 5, 3, 1, 6, 3, 1, 7, 3, 1, 8, 3, 1,
	2, 4, 1, 3, 4, 1, 4, 4, 1, 5, 4, 1,
	3, 5, 1, 0, 6, 1, 1, 6, 1, 2, 6, 1,
}

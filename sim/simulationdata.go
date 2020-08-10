package sim

import "sort"

// simulationData defines all cell data for a simulation.
type simulationData struct {
	// cellData defines the contents of each cell.
	// Each entry defines the X and Y coordinates of a cell,
	// along with its state.
	cellData CellList

	// tempData is kept in sync with cellData as far as size
	// goes. it is used as a temporary buffer in the step() function.
	tempData CellList

	// neighbours contains one entry for each cell. Each entry
	// defines the indices of all 8 neighbouring cells, or -1
	// if no neighbour exists at a specific place.
	neighbours []int

	// cellsChanged signals to a caller that the cell buffer has changed.
	cellsChanged bool

	// Instructs simulator to recompute cell neighbours before performing
	// a step call. This needs to be done when cells have been altered
	// through Set().
	staleNeighbours bool
}

// CellCount returns the number of cells in the simulation.
func (s *simulationData) CellCount() int {
	return len(s.cellData) / 3
}

// update recreates/resizes the neighbours and tempdata sets where needed.
func (s *simulationData) update() {
	cc := s.cellData.Len()

	if n := cc * 8; len(s.neighbours) > n {
		s.neighbours = s.neighbours[:n]
	} else {
		s.neighbours = make([]int, n)
	}

	if n := cc * 3; len(s.tempData) > n {
		s.tempData = s.tempData[:n]
	} else {
		s.tempData = make(CellList, n)
	}

	s.cellsChanged = true
	s.staleNeighbours = true
}

// Sort sorts the cell list.
func (s *simulationData) Sort() {
	sort.Sort(s)
}

// Load loads c2 to c1 and returns the resulting set.
// c2's top-right corner is placed at the given position.
// This ensures no cell duplicates are added.
func (s *simulationData) Load(x, y int32, v CellList) {
	s.cellData = s.cellData.Load(x, y, v)
	s.update()
}

// Unload removes all cells in v from the simulation.
func (s *simulationData) Unload(v CellList) {
	s.cellData = s.cellData.Unload(v)
	s.update()
}

// UpdateList overwrites the values of cells from v in the simulation.
func (s *simulationData) UpdateList(v CellList) {
	if len(v) == 0 {
		return
	}

	cd := s.cellData

	for i := 0; i < len(v)-2; i += 3 {
		n := s.cellData.IndexOf(v[i], v[i+1])
		if n > -1 {
			cd[n+2] = v[i+2]
		}
	}

	s.cellsChanged = true
}

func (s *simulationData) Set(x, y, state int32) {
	n := s.cellData.IndexOf(x, y)
	if n > -1 {
		s.cellsChanged = s.cellData[n+2] != state
		s.cellData[n+2] = state
		return
	}

	// Cell is new. If the new state is CellEmpty, just ignore it.
	if state == CellEmpty {
		return
	}

	s.cellData = append(s.cellData, x, y, state)
	s.tempData = append(s.tempData, 0, 0, 0)
	s.neighbours = append(s.neighbours, 0, 0, 0, 0, 0, 0, 0, 0)
	s.cellsChanged = true
	s.staleNeighbours = true
}

// Step performs a single simulation step by applying the Wireworld rules to the cell data.
func (s *simulationData) Step() {
	// Recompute neighbours if necessary.
	if s.staleNeighbours {
		s.computeNeighbours()
		s.staleNeighbours = false
	}

	var i, j, k, ci, ni int
	var x, y, v int32
	var n []int

	t0 := data.cellData
	t1 := data.tempData
	cn := data.neighbours

	for i = 0; i < len(t0)/3; i++ {
		ci, ni = i*3, i*8
		x = t0[ci+0]
		y = t0[ci+1]
		v = t0[ci+2]
		n = cn[ni : ni+8]

		switch v {
		case CellWire:
			k = 0

			for _, j = range n {
				if j > -1 && t0[j+2] == CellHead {
					k++
				}
			}

			if k == 1 || k == 2 {
				v = CellHead
			}
		case CellHead:
			v = CellTail
		case CellTail:
			v = CellWire
		}

		t1[ci] = x
		t1[ci+1] = y
		t1[ci+2] = v
	}

	// Swap buffers to make new celldata the current set.
	s.cellData = t1
	s.tempData = t0

	s.cellsChanged = true
}

// computeNeighbours recomputes all neighbours for all cells.
func (s *simulationData) computeNeighbours() {
	var cn []int
	cells := s.cellData

	for i := 0; i < len(cells)/3; i++ {
		ci, ni := i*3, i*8
		x, y := cells[ci], cells[ci+1]

		cn = s.neighbours[ni:]

		// N, S, W, E neighbours
		cn[0] = cells.IndexOf(x, y-1)
		cn[1] = cells.IndexOf(x, y+1)
		cn[2] = cells.IndexOf(x-1, y)
		cn[3] = cells.IndexOf(x+1, y)

		// Diagonal neighbours.
		cn[4] = cells.IndexOf(x-1, y-1)
		cn[5] = cells.IndexOf(x+1, y-1)
		cn[6] = cells.IndexOf(x-1, y+1)
		cn[7] = cells.IndexOf(x+1, y+1)
	}
}

func (s *simulationData) Trim() {
	s.cellData = s.cellData.Trim()
	s.update()
}

// sort interface implementation.
func (s *simulationData) Len() int {
	return s.cellData.Len()
}

// sort interface implementation.
func (s *simulationData) Less(i, j int) bool {
	return s.cellData.Less(i, j)
}

// sort interface implementation.
// This swaps elements in the celldata and neighbours lists.
func (s *simulationData) Swap(i, j int) {
	s.cellData.Swap(i, j)

	var neighbours [8]int
	nd := s.neighbours
	ni, nj := i*8, j*8

	copy(neighbours[:], nd[ni:])
	copy(nd[ni:], nd[nj:nj+8])
	copy(nd[nj:], neighbours[:])
}

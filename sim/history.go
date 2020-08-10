package sim

type historyState struct {
	before CellList
	after  CellList
}

// History provides undo/redo facilities for a simulation.
//
// It maintains undo and redo stacks for all user-edited chunks
// of cells in a simulation.
type History struct {
	undo []historyState
	redo []historyState
}

// changeHandler is called when the given cells have changed.
func (h *History) changeHandler(before, after CellList) {
	h.undo = append(h.undo, historyState{before, after})
	h.redo = h.redo[:0]
}

// Undo undoes the last cell change.
func (h *History) Undo() {
	if len(h.undo) == 0 {
		return
	}

	state := h.undo[len(h.undo)-1]
	h.undo = h.undo[:len(h.undo)-1]
	h.redo = append(h.redo, state)

	SetList(state.before)
}

// Redo redoes the last cell change.
func (h *History) Redo() {
	if len(h.redo) == 0 {
		return
	}

	state := h.redo[len(h.redo)-1]
	h.redo = h.redo[:len(h.redo)-1]
	h.undo = append(h.undo, state)

	SetList(state.after)
}

package backtrack

import (
	"slices"
	"sudoku-solver/sudoku"
)

type cellState struct {
	HasValue      bool
	Value         int
	Possibilities []int
}

type cellsState map[sudoku.Coordinate]cellState

var _ sudoku.Solution = cellsState(nil)

func (s cellState) Copy() cellState {
	newPossibilities := make([]int, len(s.Possibilities))
	copy(newPossibilities, s.Possibilities)
	return cellState{
		HasValue:      s.HasValue,
		Value:         s.Value,
		Possibilities: newPossibilities,
	}
}

func FixedCellState(value int) cellState {
	return cellState{
		HasValue:      true,
		Value:         value,
		Possibilities: nil,
	}
}

func VariableCellState(possibilities []int) cellState {
	return cellState{
		HasValue:      false,
		Value:         0,
		Possibilities: possibilities,
	}
}

// WithRemovedPossibilities removes the given values from the possibilities of this cell.
// It returns a new cellState and a bool that is true if this cell is still solvable.
// If the cell has a value, this method just returns the cellState and true.
// If the cell has no value the method returns a new cellState that has the given values
// removed from its possibilities and returns true if the new cellState has at least one
// possibility left.
func (s cellState) WithRemovedPossibilities(values ...int) (cellState, bool) {
	if s.HasValue {
		return s, true
	}
	newPossibilities := Filter(s.Possibilities, func(value int) bool {
		return !slices.Contains(values, value)
	})
	return VariableCellState(newPossibilities), len(newPossibilities) > 0
}

// WithConstrainedPossibilities returns a new cellState that only allows the given values.
// It returns the new cellState and a bool that is true if this cell is still solvable.
// If the cell has a value, this method returns true if the value is in the given values.
// If the cell has no value, this method returns a new cellState that only allows the
// given values and returns true if the new cellState has at least one possibility left.
func (s cellState) WithConstrainedPossibilities(values ...int) (cellState, bool) {
	if s.HasValue {
		return s, slices.Contains(values, s.Value)
	}
	newPossibilities := make([]int, 0, len(s.Possibilities))
	for _, value := range s.Possibilities {
		if slices.Contains(values, value) {
			newPossibilities = append(newPossibilities, value)
		}
	}
	return VariableCellState(newPossibilities), len(newPossibilities) > 0
}

func (s cellsState) Get(coordinate sudoku.Coordinate) (int, bool) {
	coorState, ok := s[coordinate]
	if !ok {
		return 0, false
	}
	if !coorState.HasValue {
		return 0, false
	}
	return coorState.Value, true
}

func (s cellsState) Copy() cellsState {
	newState := make(cellsState, len(s))
	for coor, coorState := range s {
		newState[coor] = coorState.Copy()
	}
	return newState
}

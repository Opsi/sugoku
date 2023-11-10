package backtrack

import "slices"

type cellState struct {
	HasValue      bool
	Value         int
	Possibilities []int
}

func (s cellState) Copy() cellState {
	newPossibilities := make([]int, len(s.Possibilities))
	copy(newPossibilities, s.Possibilities)
	return cellState{
		HasValue:      s.HasValue,
		Value:         s.Value,
		Possibilities: newPossibilities,
	}
}

// RemovePossibilities removes the given values from the possibilities of this cell.
// If the cell has a value, this method does nothing and returns true.
// If the cell has no value and the given values this method removes all given values
// from the possibilities and return true if the cell has at least one possibility left.
func (s *cellState) RemovePossibilities(values ...int) bool {
	if s.HasValue {
		return true
	}
	newPossibilities := make([]int, 0, len(s.Possibilities))
	for _, value := range s.Possibilities {
		if !slices.Contains(values, value) {
			newPossibilities = append(newPossibilities, value)
		}
	}
	s.Possibilities = newPossibilities
	return len(s.Possibilities) > 0
}

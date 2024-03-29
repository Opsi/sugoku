package constraint

import "sudoku-solver/sudoku"

type FixedValueConstraint struct {
	Coordinate sudoku.Coordinate
	Value      int
}

var _ sudoku.Constraint = FixedValueConstraint{}

func (c FixedValueConstraint) IsViolated(solution sudoku.Solution) bool {
	value, ok := solution.Get(c.Coordinate)
	if !ok {
		return false
	}
	return value != c.Value
}

func (c FixedValueConstraint) ConstrainedCoordinates() []sudoku.Coordinate {
	return []sudoku.Coordinate{c.Coordinate}
}

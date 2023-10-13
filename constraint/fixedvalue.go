package constraint

import "sudoku-solver/sudoku"

type FixedValueConstraint struct {
	Coordinate sudoku.Coordinate
	Value      int
}

var _ sudoku.Constraint = FixedValueConstraint{}

func (c FixedValueConstraint) Check(solution sudoku.Solution) sudoku.ConstraintResult {
	value, ok := solution[c.Coordinate]
	if !ok {
		return sudoku.ConstraintResultValidAndNotSolved
	}
	if value != c.Value {
		return sudoku.ConstraintResultInvalid
	}
	return sudoku.ConstraintResultValidAndSolved
}

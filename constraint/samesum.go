package constraint

import (
	"fmt"
	"sudoku-solver/sudoku"
)

// SameSumConstraint is a constraint that requires that the sum of the values
// at the coordinates in Coordinates1 is the same as the sum of the values at
// the coordinates in Coordinates2.
type SameSumConstraint struct {
	Coordinates1 []sudoku.Coordinate
	Coordinates2 []sudoku.Coordinate
}

var _ sudoku.Constraint = SameSumConstraint{}

func (c SameSumConstraint) IsViolated(solution sudoku.Solution) bool {
	sum1 := 0
	for _, coor := range c.Coordinates1 {
		value, ok := solution[coor]
		if !ok {
			return false
		}
		sum1 += value
	}
	sum2 := 0
	for _, coor := range c.Coordinates2 {
		value, ok := solution[coor]
		if !ok {
			return false
		}
		sum2 += value
	}
	return sum1 != sum2
}

// NewArrowConstrain creates a new SameSumConstraint that requires that the sum
// of the values at the coordinates in path is the same as the value at circle.
func NewArrowConstrain(circle sudoku.Coordinate, path []sudoku.Coordinate) (*SameSumConstraint, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("arrow path must not be empty")
	}
	return &SameSumConstraint{
		Coordinates1: []sudoku.Coordinate{circle},
		Coordinates2: path,
	}, nil
}

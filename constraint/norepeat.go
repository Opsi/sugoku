package constraint

import (
	"fmt"
	"sudoku-solver/sudoku"
)

type NoRepeatConstraint struct {
	Coordinates []sudoku.Coordinate
}

var _ sudoku.Constraint = NoRepeatConstraint{}

func (c NoRepeatConstraint) Check(solution sudoku.Solution) sudoku.ConstraintResult {
	seen := make(map[int]struct{})
	for _, cell := range c.Coordinates {
		value, ok := solution[cell]
		if !ok {
			return sudoku.ConstraintResultValidAndNotSolved
		}
		if _, ok := seen[value]; ok {
			return sudoku.ConstraintResultInvalid
		}
		seen[value] = struct{}{}
	}
	return sudoku.ConstraintResultValidAndSolved
}

func NewRowConstraint(row int, coordinates []sudoku.Coordinate) (NoRepeatConstraint, error) {
	rowCoords := make([]sudoku.Coordinate, 0)
	for _, coordinate := range coordinates {
		if coordinate.Row == row {
			rowCoords = append(rowCoords, coordinate)
		}
	}
	if len(rowCoords) == 0 {
		return NoRepeatConstraint{}, fmt.Errorf("no coordinates found for row %d", row)
	}
	return NoRepeatConstraint{Coordinates: rowCoords}, nil
}

func NewColumnConstraint(col int, coordinates []sudoku.Coordinate) (NoRepeatConstraint, error) {
	colCoords := make([]sudoku.Coordinate, 0)
	for _, coordinate := range coordinates {
		if coordinate.Col == col {
			colCoords = append(colCoords, coordinate)
		}
	}
	if len(colCoords) == 0 {
		return NoRepeatConstraint{}, fmt.Errorf("no coordinates found for column %d", col)
	}
	return NoRepeatConstraint{Coordinates: colCoords}, nil
}

func NewSquareConstraint(row, col, length int, coordinates []sudoku.Coordinate) (NoRepeatConstraint, error) {
	squareCoords := make([]sudoku.Coordinate, 0)
	for _, coordinate := range coordinates {
		if coordinate.Row < row {
			continue
		}
		if coordinate.Row >= row+length {
			continue
		}
		if coordinate.Col < col {
			continue
		}
		if coordinate.Col >= col+length {
			continue
		}
		squareCoords = append(squareCoords, coordinate)
	}
	if len(squareCoords) == 0 {
		return NoRepeatConstraint{}, fmt.Errorf("no coordinates found for square of length %d at row %d and column %d", length, row, col)
	}
	return NoRepeatConstraint{Coordinates: squareCoords}, nil
}

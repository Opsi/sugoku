package constraint

import (
	"fmt"
	"math"
	"sudoku-solver/sudoku"
)

type NoRepeatConstraint struct {
	Coordinates []sudoku.Coordinate
}

var _ sudoku.Constraint = NoRepeatConstraint{}

func (c NoRepeatConstraint) IsViolated(solution sudoku.Solution) bool {
	seen := make(map[int]struct{})
	for _, coord := range c.Coordinates {
		value, ok := solution.Get(coord)
		if !ok {
			continue
		}
		if _, ok := seen[value]; ok {
			return true
		}
		seen[value] = struct{}{}
	}
	return false
}

func (c NoRepeatConstraint) ConstrainedCoordinates() []sudoku.Coordinate {
	return c.Coordinates
}

func RowConstraint(row int, coordinates []sudoku.Coordinate) (*NoRepeatConstraint, error) {
	rowCoords := make([]sudoku.Coordinate, 0)
	for _, coordinate := range coordinates {
		if coordinate.Row == row {
			rowCoords = append(rowCoords, coordinate)
		}
	}
	if len(rowCoords) == 0 {
		return nil, fmt.Errorf("no coordinates found for row %d", row)
	}
	return &NoRepeatConstraint{Coordinates: rowCoords}, nil
}

func RowConstraints(coordinates []sudoku.Coordinate) ([]NoRepeatConstraint, error) {
	constraints := make([]NoRepeatConstraint, 0)
	seenRows := make(map[int]struct{})
	for _, coordinate := range coordinates {
		if _, ok := seenRows[coordinate.Row]; ok {
			continue
		}
		constraint, err := RowConstraint(coordinate.Row, coordinates)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, *constraint)
		seenRows[coordinate.Row] = struct{}{}
	}
	return constraints, nil
}

func ColumnConstraint(col int, coordinates []sudoku.Coordinate) (*NoRepeatConstraint, error) {
	colCoords := make([]sudoku.Coordinate, 0)
	for _, coordinate := range coordinates {
		if coordinate.Col == col {
			colCoords = append(colCoords, coordinate)
		}
	}
	if len(colCoords) == 0 {
		return nil, fmt.Errorf("no coordinates found for column %d", col)
	}
	return &NoRepeatConstraint{Coordinates: colCoords}, nil
}

func ColumnConstraints(coordinates []sudoku.Coordinate) ([]NoRepeatConstraint, error) {
	constraints := make([]NoRepeatConstraint, 0)
	seenCols := make(map[int]struct{})
	for _, coordinate := range coordinates {
		if _, ok := seenCols[coordinate.Col]; ok {
			continue
		}
		constraint, err := ColumnConstraint(coordinate.Col, coordinates)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, *constraint)
		seenCols[coordinate.Col] = struct{}{}
	}
	return constraints, nil
}

func SquareConstraint(row, col, length int, coordinates []sudoku.Coordinate) (*NoRepeatConstraint, error) {
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
		return nil, fmt.Errorf("no coordinates found for square of length %d at row %d and column %d", length, row, col)
	}
	return &NoRepeatConstraint{Coordinates: squareCoords}, nil
}

func BoxConstraints(coordinates []sudoku.Coordinate) ([]NoRepeatConstraint, error) {
	if len(coordinates) == 0 {
		return nil, fmt.Errorf("no coordinates given")
	}
	// first find the bounds of the coordinates
	minRow := math.MaxInt
	maxRow := math.MinInt
	minCol := math.MaxInt
	maxCol := math.MinInt
	for _, coordinate := range coordinates {
		if coordinate.Row < minRow {
			minRow = coordinate.Row
		}
		if coordinate.Row > maxRow {
			maxRow = coordinate.Row
		}
		if coordinate.Col < minCol {
			minCol = coordinate.Col
		}
		if coordinate.Col > maxCol {
			maxCol = coordinate.Col
		}
	}

	rows := maxRow - minRow + 1
	cols := maxCol - minCol + 1
	if rows != cols {
		return nil, fmt.Errorf("coordinates are not square: %d rows and %d columns", rows, cols)
	}
	boxSize := int(math.Sqrt(float64(rows)))
	if rows != boxSize*boxSize {
		return nil, fmt.Errorf("field size %d is not a square number", rows)
	}

	constraints := make([]NoRepeatConstraint, 0)
	for row := minRow; row <= maxRow; row += boxSize {
		for col := minCol; col <= maxCol; col += boxSize {
			constraint, err := SquareConstraint(row, col, boxSize, coordinates)
			if err != nil {
				return nil, fmt.Errorf("create box constraint for row %d and column %d: %w", row, col, err)
			}
			constraints = append(constraints, *constraint)
		}
	}
	return constraints, nil
}

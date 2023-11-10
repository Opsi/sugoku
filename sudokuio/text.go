package sudokuio

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

func ParseStringRow(row string) ([]int, error) {
	values := make([]int, 0)
	for charIndex, char := range row {
		switch char {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			values = append(values, int(char-'0'))
		case '-':
			values = append(values, 0)
		case ' ':
			// ignore
			continue
		default:
			return nil, fmt.Errorf("%d. character '%c' is invalid", charIndex+1, char)
		}
	}
	return values, nil
}

func ParseStringRows(input string) ([][]int, error) {
	rows := make([][]int, 0)
	expectedRowLength := -1
	for lineIndex, line := range strings.Split(input, "\n") {
		row, err := ParseStringRow(line)
		if err != nil {
			return nil, fmt.Errorf("parse %d. row: %w", lineIndex+1, err)
		}
		if expectedRowLength == -1 {
			expectedRowLength = len(row)
		}
		if len(row) != expectedRowLength {
			return nil, fmt.Errorf("parse %d. row: expected %d cells, got %d", lineIndex+1, expectedRowLength, len(row))
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func StringRowsToSudoku(rowStrings []string) (*sudoku.Sudoku, error) {
	// FIXME: this is bit sneaky
	input := strings.Join(rowStrings, "\n")
	rows, err := ParseStringRows(input)
	if err != nil {
		return nil, fmt.Errorf("parse rows: %w", err)
	}

	// TODO: support different sudoku sizes
	if len(rows) != 9 {
		return nil, fmt.Errorf("expected 9 rows, got %d", len(rows))
	}
	if len(rows[0]) != 9 {
		return nil, fmt.Errorf("expected 9 columns, got %d", len(rows[0]))
	}
	possibleValues := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	constraints := make([]sudoku.Constraint, 0)
	coordinates := make([]sudoku.Coordinate, 0, 81)
	for rowIdx, row := range rows {
		for colIdx, value := range row {
			coordinate := sudoku.Coordinate{
				Row: rowIdx + 1,
				Col: colIdx + 1,
			}
			coordinates = append(coordinates, coordinate)
			if value == 0 {
				continue
			}
			if !slices.Contains(possibleValues, value) {
				return nil, fmt.Errorf("invalid value %d in row %d, column %d", value, rowIdx+1, colIdx+1)
			}
			constraint := constraint.FixedValueConstraint{
				Coordinate: coordinate,
				Value:      value,
			}
			constraints = append(constraints, constraint)
		}
	}

	return &sudoku.Sudoku{
		Coordinates:    coordinates,
		PossibleValues: possibleValues,
		Constraints:    constraints,
	}, nil
}

func ParseString(input string) (*sudoku.Sudoku, error) {
	tokens := make([]rune, 0, len(input))
	// tokenize check if string only contains valid characters: 1-9, -, \n, \r, \t, space
	for _, c := range input {
		switch c {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			tokens = append(tokens, c)
		case '\n', '\r', '\t', ' ':
			// ignore
		default:
			return nil, fmt.Errorf("invalid character: %c", c)
		}
	}

	// check if the number of tokens is too small or too big
	if len(tokens) < 1 {
		return nil, fmt.Errorf("sudoku is empty")
	}

	if len(tokens) > 81 {
		return nil, fmt.Errorf("sudoku is too large with %d cells", len(tokens))
	}

	// check if the number of tokens builds a perfect square
	squareRoot := math.Sqrt(float64(len(tokens)))
	if squareRoot != math.Floor(squareRoot) {
		return nil, fmt.Errorf("amount of tokens %d is not a square number", len(tokens))
	}
	sodokuLength := int(squareRoot)

	// check if the sodoku size is a square number (to determine the box size)
	squareRoot = math.Sqrt(float64(sodokuLength))
	if squareRoot != math.Floor(squareRoot) {
		return nil, fmt.Errorf("amount of tokens %d is not a cube number", len(tokens))
	}

	possibleValues := make([]int, 0, sodokuLength)
	for i := 1; i <= sodokuLength; i++ {
		possibleValues = append(possibleValues, i)
	}

	coordinates := make([]sudoku.Coordinate, 0, len(tokens))
	for x := 1; x <= sodokuLength; x++ {
		for y := 1; y <= sodokuLength; y++ {
			coordinates = append(coordinates, sudoku.Coordinate{Row: x, Col: y})
		}
	}

	constraints := make([]sudoku.Constraint, 0)

	// add row constraints
	rowConstraints, err := constraint.RowConstraints(coordinates)
	if err != nil {
		// this should never happen
		panic(err)
	}
	for _, rowConstraint := range rowConstraints {
		constraints = append(constraints, rowConstraint)
	}

	// add column constraints
	colConstraints, err := constraint.ColumnConstraints(coordinates)
	if err != nil {
		// this should never happen
		panic(err)
	}
	for _, colConstraint := range colConstraints {
		constraints = append(constraints, colConstraint)
	}

	// add box constraints
	boxConstrains, err := constraint.BoxConstraints(coordinates)
	if err != nil {
		// this should never happen
		panic(err)
	}
	for _, boxConstrain := range boxConstrains {
		constraints = append(constraints, boxConstrain)
	}

	// add fixed value constraints
	for i, token := range tokens {
		if token == '-' {
			continue
		}
		coordinate := coordinates[i]
		value := int(token - '0')
		fixedValueConstraint := constraint.FixedValueConstraint{
			Coordinate: coordinate,
			Value:      value,
		}
		constraints = append(constraints, fixedValueConstraint)
	}

	// create an array of only the valid characters: 1-9 and -
	return &sudoku.Sudoku{
		Coordinates:    coordinates,
		PossibleValues: possibleValues,
		Constraints:    constraints,
	}, nil
}

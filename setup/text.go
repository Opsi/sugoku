package setup

import (
	"fmt"
	"math"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

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
	boxLength := int(squareRoot)

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
	for x := 1; x <= sodokuLength; x++ {
		rowConstraint, err := constraint.NewRowConstraint(x, coordinates)
		if err != nil {
			// this should never happen
			panic(err)
		}
		constraints = append(constraints, rowConstraint)
	}

	// add column constraints
	for y := 1; y <= sodokuLength; y++ {
		colConstraint, err := constraint.NewColumnConstraint(y, coordinates)
		if err != nil {
			// this should never happen
			panic(err)
		}
		constraints = append(constraints, colConstraint)
	}

	// add box constraints
	for x := 1; x <= sodokuLength; x += boxLength {
		for y := 1; y <= sodokuLength; y += boxLength {
			boxConstraint, err := constraint.NewSquareConstraint(x, y, boxLength, coordinates)
			if err != nil {
				// this should never happen
				panic(err)
			}
			constraints = append(constraints, boxConstraint)
		}
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

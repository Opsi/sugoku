package sudokuio

import (
	"encoding/json"
	"fmt"
	"slices"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

// RawConstraintGen is a type that can be unmarshalled from a JSON object.
// It represents a function that takes a sudoku and returns a slice of constraints.
type RawConstraintGen func(sudoku.Sudoku) ([]sudoku.Constraint, error)

var _ json.Unmarshaler = (*RawConstraintGen)(nil)

type constraintType string

const (
	constraintTypeNormalSudokuRules constraintType = "normalSudokuRules"
	constraintTypeArrow             constraintType = "arrow"
	constraintTypeFixedValues       constraintType = "fixedValues"
)

type baseConstraintGen struct {
	Type constraintType `json:"type"`
}

func generateNormalSudokuRules(sudoku.Sudoku) ([]sudoku.Constraint, error) {
	// TODO: implement
	return []sudoku.Constraint{}, nil
}

type arrowConstraintGen struct {
	Circle RawCoordinate   `json:"circle"`
	Path   []RawCoordinate `json:"path"`
}

func (g arrowConstraintGen) generate(s sudoku.Sudoku) ([]sudoku.Constraint, error) {
	circle := sudoku.Coordinate(g.Circle)
	// check that the circle is in the sudoku
	if ok := slices.Contains(s.Coordinates, circle); ok {
		return nil, fmt.Errorf("arrow circle coordinate %s is not in the sudoku", circle)
	}
	path := make([]sudoku.Coordinate, len(g.Path))
	for _, c := range g.Path {
		coord := sudoku.Coordinate(c)
		if ok := slices.Contains(s.Coordinates, coord); ok {
			return nil, fmt.Errorf("arrow path coordinate %s is not in the sudoku", coord)
		}
		path = append(path, coord)
	}
	arrow, err := constraint.NewArrowConstraint(circle, path)
	if err != nil {
		return nil, fmt.Errorf("invalid arrow constraint: %w", err)
	}
	return []sudoku.Constraint{arrow}, nil
}

type fixedValuesConstraintGen struct {
	Values map[RawCoordinate]int `json:"values"`
}

func (g fixedValuesConstraintGen) generate(s sudoku.Sudoku) ([]sudoku.Constraint, error) {
	constraints := make([]sudoku.Constraint, 0, len(g.Values))
	for c, v := range g.Values {
		coord := sudoku.Coordinate(c)
		if ok := slices.Contains(s.Coordinates, coord); ok {
			return nil, fmt.Errorf("fixed value coordinate %s is not in the sudoku", coord)
		}
		if ok := slices.Contains(s.PossibleValues, v); ok {
			return nil, fmt.Errorf("fixed value %d is not allowed in the sudoku", v)
		}
		constraints = append(constraints, constraint.FixedValueConstraint{
			Coordinate: coord,
			Value:      v,
		})
	}
	if len(constraints) == 0 {
		return nil, fmt.Errorf("fixed values constraint must have at least one value")
	}
	return constraints, nil
}

func (c *RawConstraintGen) UnmarshalJSON(data []byte) error {
	var base baseConstraintGen
	if err := json.Unmarshal(data, &base); err != nil {
		return fmt.Errorf("constraint should be an object, got %s", data)
	}
	switch base.Type {
	case constraintTypeNormalSudokuRules:
		*c = generateNormalSudokuRules
		return nil
	case constraintTypeArrow:
		var arrowGen arrowConstraintGen
		if err := json.Unmarshal(data, &arrowGen); err != nil {
			return fmt.Errorf("invalid arrow constraint: %w", err)
		}
		*c = arrowGen.generate
		return nil
	case constraintTypeFixedValues:
		var fixedValuesGen fixedValuesConstraintGen
		if err := json.Unmarshal(data, &fixedValuesGen); err != nil {
			return fmt.Errorf("invalid fixed values constraint: %w", err)
		}
		*c = fixedValuesGen.generate
		return nil
	default:
		return fmt.Errorf("unknown constraint type %s", base.Type)
	}
}

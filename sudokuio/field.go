package sudokuio

import (
	"encoding/json"
	"fmt"
	"sudoku-solver/sudoku"
)

type RawFieldGen func() (sudoku.Sudoku, error)

var _ json.Unmarshaler = (*RawFieldGen)(nil)

type fieldType string

const (
	fieldTypeNormal fieldType = "normal"
)

type baseFieldGen struct {
	Type fieldType `json:"type"`
}

func (g *RawFieldGen) UnmarshalJSON(data []byte) error {
	var base baseFieldGen
	if err := json.Unmarshal(data, &base); err != nil {
		return fmt.Errorf("field should be an object: %w", err)
	}
	switch base.Type {
	case fieldTypeNormal:
		*g = generateNormalSudoku
		return nil
	default:
		return fmt.Errorf("unknown sudoku type %q", base.Type)
	}
}

func generateNormalSudoku() (sudoku.Sudoku, error) {
	coordinates := make([]sudoku.Coordinate, 0, 81)
	for row := 1; row <= 9; row++ {
		for col := 1; col <= 9; col++ {
			coordinates = append(coordinates, sudoku.Coordinate{
				Row: row,
				Col: col,
			})
		}
	}
	possibleValues := make([]int, 0, 9)
	for i := 1; i <= 9; i++ {
		possibleValues = append(possibleValues, i)
	}
	return sudoku.Sudoku{
		Coordinates:    coordinates,
		PossibleValues: possibleValues,
	}, nil
}

package sudokuio

import (
	"encoding/json"
	"fmt"
	"sudoku-solver/sudoku"
)

type RawFieldGen func() (*sudoku.Sudoku, error)

var _ json.Unmarshaler = (*RawFieldGen)(nil)

type fieldType string

const (
	fieldTypeNormal fieldType = "normal"
)

type baseFieldGen struct {
	Type fieldType `json:"type"`
}

type normalFieldGen struct {
	Rows []string
}

func (g *RawFieldGen) UnmarshalJSON(data []byte) error {
	var base baseFieldGen
	if err := json.Unmarshal(data, &base); err != nil {
		return fmt.Errorf("field should be an object: %w", err)
	}
	switch base.Type {
	case fieldTypeNormal:
		var normalGen normalFieldGen
		if err := json.Unmarshal(data, &normalGen); err != nil {
			return fmt.Errorf("parse normal field: %w", err)
		}
		*g = normalGen.generate
		return nil
	default:
		return fmt.Errorf("unknown sudoku type %q", base.Type)
	}
}

func (n normalFieldGen) generate() (*sudoku.Sudoku, error) {
	return StringRowsToSudoku(n.Rows)
}

package sudokuio

import (
	"encoding/json"
	"fmt"
	"sudoku-solver/sudoku"
)

type rawJsonSudoku struct {
	Field       RawFieldGen        `json:"field"`
	Constraints []RawConstraintGen `json:"constraints"`
}

func (s rawJsonSudoku) generate() (*sudoku.Sudoku, error) {
	if s.Field == nil {
		return nil, fmt.Errorf("field is missing")
	}
	sudok, err := s.Field()
	if err != nil {
		return nil, fmt.Errorf("generate field: %w", err)
	}
	for _, constraintGen := range s.Constraints {
		constraints, err := constraintGen(sudok)
		if err != nil {
			return nil, fmt.Errorf("generate constraints: %w", err)
		}
		for ci, c := range constraints {
			fmt.Printf("[%d] constraint: %+v\n", ci, c)
		}

		sudok.Constraints = append(sudok.Constraints, constraints...)
	}
	return &sudok, nil
}

func ParseJSON(bytes []byte) (*sudoku.Sudoku, error) {
	var raw rawJsonSudoku
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	sudok, err := raw.generate()
	if err != nil {
		return nil, fmt.Errorf("generate sudoku: %w", err)
	}
	return sudok, nil
}

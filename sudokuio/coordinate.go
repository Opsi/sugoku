package sudokuio

import (
	"encoding/json"
	"fmt"
	"sudoku-solver/sudoku"
)

// RawCoordinate is a type that can be unmarshalled from a JSON string of the form "R1C1" or "C9R3"
type RawCoordinate sudoku.Coordinate

var _ json.Unmarshaler = (*RawCoordinate)(nil)

func (c *RawCoordinate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("coordinate should be a string, got %s", data)
	}
	parsed, err := sudoku.ParseCoordinateString(s)
	if err != nil {
		return fmt.Errorf("parse coordinate string: %w", err)
	}
	*c = RawCoordinate(parsed)
	return nil
}

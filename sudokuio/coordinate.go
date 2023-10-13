package sudokuio

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"sudoku-solver/sudoku"
)

// RawCoordinate is a type that can be unmarshalled from a JSON string of the form "R1C1" or "C9R3"
type RawCoordinate sudoku.Coordinate

var errWrongFromat = fmt.Errorf("coordinate should be of the form R1C4 or C7R2")

var _ json.Unmarshaler = (*RawCoordinate)(nil)

var coordinateRegex = regexp.MustCompile(`^([RC])(\d+)([RC])(\d+)$`)

func (c *RawCoordinate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("coordinate should be a string, got %s", data)
	}
	// check that the coordinate matches the regex
	if ok := coordinateRegex.MatchString(s); !ok {
		return fmt.Errorf("%w: %s", errWrongFromat, s)
	}

	// extract the row and column numbers from the string
	matches := coordinateRegex.FindStringSubmatch(s)

	if len(matches) != 5 {
		return fmt.Errorf("invalid coordinate: %s", s)
	}

	if matches[1] == matches[3] {
		// wrote R1R2 or C6C5
		return fmt.Errorf("%w: %s", errWrongFromat, s)
	}

	if matches[1] != "R" && matches[1] != "C" {
		return fmt.Errorf("invalid row/column specifier: %s", matches[1])
	}

	if matches[3] != "R" && matches[3] != "C" {
		return fmt.Errorf("invalid row/column specifier: %s", matches[3])
	}

	rowIndex := 2
	colIndex := 4
	if matches[1] == "C" {
		rowIndex, colIndex = colIndex, rowIndex
	}

	rowStr, colStr := matches[rowIndex], matches[colIndex]
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return fmt.Errorf("invalid row number: %s", rowStr)
	}
	col, err := strconv.Atoi(colStr)
	if err != nil {
		return fmt.Errorf("invalid column number: %s", colStr)
	}

	*c = RawCoordinate(sudoku.Coordinate{Row: row, Col: col})
	return nil
}

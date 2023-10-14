package sudoku

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var coordinateRegex = regexp.MustCompile(`^([RC])(\d+)([RC])(\d+)$`)

var errWrongFromat = fmt.Errorf("coordinate should be of the form R1C4 or C7R2")

type Coordinate struct {
	Row, Col int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("R%d-C%d", c.Row, c.Col)
}

func ParseCoordinateString(s string) (Coordinate, error) {
	s = strings.ToUpper(s)
	// check that the coordinate matches the regex
	if ok := coordinateRegex.MatchString(s); !ok {
		return Coordinate{}, fmt.Errorf("%w: %s", errWrongFromat, s)
	}

	// extract the row and column numbers from the string
	matches := coordinateRegex.FindStringSubmatch(s)

	if len(matches) != 5 {
		return Coordinate{}, fmt.Errorf("invalid coordinate: %s", s)
	}

	if matches[1] == matches[3] {
		// wrote R1R2 or C6C5
		return Coordinate{}, fmt.Errorf("%w: %s", errWrongFromat, s)
	}

	if matches[1] != "R" && matches[1] != "C" {
		return Coordinate{}, fmt.Errorf("invalid row/column specifier: %s", matches[1])
	}

	if matches[3] != "R" && matches[3] != "C" {
		return Coordinate{}, fmt.Errorf("invalid row/column specifier: %s", matches[3])
	}

	rowIndex := 2
	colIndex := 4
	if matches[1] == "C" {
		rowIndex, colIndex = colIndex, rowIndex
	}

	rowStr, colStr := matches[rowIndex], matches[colIndex]
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid row number: %s", rowStr)
	}
	col, err := strconv.Atoi(colStr)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid column number: %s", colStr)
	}

	return Coordinate{Row: row, Col: col}, nil
}

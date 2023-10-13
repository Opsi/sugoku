package sudokuio_test

import (
	"sudoku-solver/sudoku"
	"sudoku-solver/sudokuio"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalCoordinate(t *testing.T) {
	tests := []struct {
		input    string
		expected sudoku.Coordinate
	}{
		{"R1C2", sudoku.Coordinate{Row: 1, Col: 2}},
		{"C1R2", sudoku.Coordinate{Row: 2, Col: 1}},
		{"C0R9", sudoku.Coordinate{Row: 9, Col: 0}},
		{"R0C9", sudoku.Coordinate{Row: 0, Col: 9}},
		{"C10R999999", sudoku.Coordinate{Row: 999999, Col: 10}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			asJSONString := []byte(`"` + test.input + `"`)
			var c sudokuio.RawCoordinate
			err := c.UnmarshalJSON(asJSONString)
			require.NoError(t, err)
			assert.Equal(t, test.expected, sudoku.Coordinate(c))
		})
	}
}

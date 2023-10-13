package setup_test

import (
	"sudoku-solver/setup"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseText(t *testing.T) {
	testSudoku := `
	24- --- -86
	--3 --- ---
	1-- --2 5--
	59- -1- --2
	--7 --- 3--
	8-- -4- -97
	--5 8-- --3
	--- --- 6--
	32- --- -19
	`
	sudoku, err := setup.ParseString(testSudoku)
	require.NoError(t, err)
	assert.Equal(t, len(sudoku.Coordinates), 81)
	assert.Equal(t, len(sudoku.PossibleValues), 9)
	// constaints are: 9 rows, 9 columns, 9 blocks and 26 already set values
	expectedConstraints := 9 + 9 + 9 + 26
	assert.Equal(t, len(sudoku.Constraints), expectedConstraints)
}

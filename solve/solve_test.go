package solve_test

import (
	"sudoku-solver/constraint"
	"sudoku-solver/solve"
	"sudoku-solver/sudoku"
	"sudoku-solver/sudokuio"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readSudokuStr(t *testing.T, s string) sudoku.Sudoku {
	t.Helper()
	sud, err := sudokuio.ParseString(s)
	if err != nil {
		t.Fatalf("failed to parse test sudoku: %v", err)
	}
	return *sud
}

func readSolutionStr(t *testing.T, s string) sudoku.Solution {
	t.Helper()
	// we are cheeky and just parse the solution as a sudoku and read the fixed value constraints
	sud, err := sudokuio.ParseString(s)
	if err != nil {
		t.Fatalf("failed to parse test solution: %v", err)
	}
	solution := make(sudoku.Solution)
	for _, constr := range sud.Constraints {
		fvc, ok := constr.(constraint.FixedValueConstraint)
		if !ok {
			continue
		}
		solution[fvc.Coordinate] = fvc.Value
	}
	return solution
}

func TestSolve(t *testing.T) {
	const sudokuStr = `
		24- --- -86
		--3 --- ---
		1-- --2 5--
		59- -1- --2
		--7 --- 3--
		8-- -4- -97
		--5 8-- --3
		--- --- 6--
		32- --- -19`

	const solutionStr = `
		249 135 786
		753 468 921
		186 972 534
		594 713 862
		617 289 345
		832 546 197
		465 891 273
		971 324 658
		328 657 419`

	sudok := readSudokuStr(t, sudokuStr)
	solution := readSolutionStr(t, solutionStr)
	require.NoError(t, sudok.Check(solution))
	solved, err := solve.Solve(sudok)
	require.NoError(t, err)
	require.NoError(t, sudok.Check(solved))
	for solutionCoordinate, solutionValue := range solution {
		solvedValue, ok := solved[solutionCoordinate]
		require.True(t, ok)
		assert.Equal(t, solutionValue, solvedValue)
	}
}

func TestArrowSudoku(t *testing.T) {
	const sudokuStr = `
	7-- -4- --3
	--- --- ---
	--9 --- 8--
	--- 6-2 ---
	--- --- ---
	--- 7-9 ---
	--5 --- 7--
	--- --- ---
	1-- -3- --2`

	const solutionStr = `
	728 945 613
	351 876 294
	469 123 875
	983 612 457
	617 354 928
	542 789 136
	835 261 749
	274 598 361
	196 437 582`

	arrowStrSlices := [][]string{
		{"r3c2", "r2c1", "r1c2", "r2c3"},
		{"r2c5", "r3c6", "r4c7"},
		{"r3c8", "r2c7", "r1c8", "r2c9"},
		{"r6c5", "r5c4", "r4c5", "r5c6"},
		{"r5c8", "r6c7", "r7c6"},
		{"r9c2", "r8c1", "r7c2", "r8c3"},
		{"r8c5", "r7c4", "r6c3", "r5c2", "r4c3", "r3c4"},
		{"r9c8", "r8c7", "r7c8", "r8c9"},
	}

	sudok := readSudokuStr(t, sudokuStr)
	for _, arrowStrSlice := range arrowStrSlices {
		coords := make([]sudoku.Coordinate, 0, len(arrowStrSlice))
		for _, coordStr := range arrowStrSlice {
			jsonStr := []byte(`"` + coordStr + `"`)
			var rawCoord sudokuio.RawCoordinate
			err := rawCoord.UnmarshalJSON(jsonStr)
			require.NoError(t, err)
			coords = append(coords, sudoku.Coordinate(rawCoord))
		}
		arrowConstraint, err := constraint.NewArrowConstraint(coords[0], coords[1:])
		require.NoError(t, err)
		sudok.Constraints = append(sudok.Constraints, arrowConstraint)
	}

	solution := readSolutionStr(t, solutionStr)
	require.NoError(t, sudok.Check(solution))
	solved, err := solve.Solve(sudok)
	require.NoError(t, err)
	require.NoError(t, sudok.Check(solved))
	for solutionCoordinate, solutionValue := range solution {
		solvedValue, ok := solved[solutionCoordinate]
		require.True(t, ok)
		assert.Equal(t, solutionValue, solvedValue)
	}
}

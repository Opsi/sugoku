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

const testSudoku = `
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

const testSolution = `
249 135 786
753 468 921
186 972 534
594 713 862
617 289 345
832 546 197
465 891 273
971 324 658
328 657 419
`

func readTestSudoku(t *testing.T) sudoku.Sudoku {
	t.Helper()
	sud, err := sudokuio.ParseString(testSudoku)
	if err != nil {
		t.Fatalf("failed to parse test sudoku: %v", err)
	}
	return *sud
}

func readTestSolution(t *testing.T) sudoku.Solution {
	t.Helper()
	// we are cheeky and just parse the solution as a sudoku and read the fixed value constraints
	sud, err := sudokuio.ParseString(testSolution)
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
	if len(solution) != 81 {
		t.Fatalf("test solution does not contain 81 values")
	}
	return solution
}

func TestSolve(t *testing.T) {
	sudok := readTestSudoku(t)
	solution := readTestSolution(t)
	solved, err := solve.Solve(sudok)
	require.NoError(t, err)
	require.True(t, sudok.IsSolved(solved))
	for solutionCoordinate, solutionValue := range solution {
		solvedValue, ok := solved[solutionCoordinate]
		require.True(t, ok)
		assert.Equal(t, solutionValue, solvedValue)
	}
}

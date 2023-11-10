package solve

import (
	"context"
	"fmt"
	"sudoku-solver/sudoku"
)

var errInvalidSolution = fmt.Errorf("invalid solution")

// Solve takes a sudoku and returns a solution, or an error if no solution
// exists for the given sudoku or if the sudoku is invalid.
func Solve(sudok sudoku.Sudoku) (sudoku.Solution, error) {
	solution, ok := BacktrackSudoku(context.TODO(), sudok)
	if !ok {
		return nil, errInvalidSolution
	}
	return solution, nil
}

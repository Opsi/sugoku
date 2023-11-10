package solve

import (
	"context"
	"sudoku-solver/sudoku"
)

type Candidate interface {
	IsBroken() bool
	IsSolved() bool
	NextCandidates() []Candidate
	Get(coordinate sudoku.Coordinate) (int, bool)
}

func Backtrack(ctx context.Context, root Candidate) <-chan sudoku.Solution {
	solutions := make(chan sudoku.Solution, 1)
	go func() {
		defer close(solutions)
		backtrack(ctx, solutions, root)
	}()
	return solutions
}

func BacktrackSudoku(ctx context.Context, sudok sudoku.Sudoku) (sudoku.Solution, bool) {
	root := rootSimple(sudok)
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	solutions := Backtrack(subCtx, root)
	select {
	case <-ctx.Done():
		return nil, false
	case solution := <-solutions:
		return solution, true
	}
}

func backtrack(ctx context.Context, solutions chan<- sudoku.Solution, candidate Candidate) {
	if ctx.Err() != nil {
		return
	}
	if candidate.IsBroken() {
		return
	}
	if candidate.IsSolved() {
		select {
		case <-ctx.Done():
		case solutions <- candidate:
		}
		return
	}
	for _, nextCandidate := range candidate.NextCandidates() {
		backtrack(ctx, solutions, nextCandidate)
	}
}

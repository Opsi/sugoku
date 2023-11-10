package backtrack

import (
	"context"
	"fmt"
	"sudoku-solver/sudoku"
)

type Candidate interface {
	IsBroken() bool
	IsSolved() bool
	NextCandidates() []Candidate
	Get(coordinate sudoku.Coordinate) (int, bool)
}

// FindSolutions returns a channel that will be closed when all solutions have been found
// or when the context is cancelled.
func FindSolutions(ctx context.Context, mode Mode, sudok sudoku.Sudoku) <-chan sudoku.Solution {
	root := mode.RootCandidate(sudok)
	solutions := make(chan sudoku.Solution, 1)
	go func() {
		defer close(solutions)
		backtrack(ctx, solutions, root)
	}()
	return solutions
}

// FindSolution returns the first solution found or an error if no solution was found or
// the context was cancelled before a solution was found.
func FindSolution(ctx context.Context, mode Mode, sudok sudoku.Sudoku) (sudoku.Solution, error) {
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	solutions := FindSolutions(subCtx, mode, sudok)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case solution, more := <-solutions:
		if !more {
			return nil, fmt.Errorf("no solution found")
		}
		return solution, nil
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

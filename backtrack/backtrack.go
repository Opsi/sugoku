package backtrack

import (
	"context"
	"fmt"
	"sudoku-solver/sudoku"
)

type Candidate interface {
	sudoku.Solution
	NextCandidates() []Candidate
}

type backtracker struct {
	sudok     sudoku.Sudoku
	solutions chan sudoku.Solution
}

// FindSolutions returns a channel that will be closed when all solutions have been found
// or when the context is cancelled.
func FindSolutions(ctx context.Context, mode Mode, sudok sudoku.Sudoku) <-chan sudoku.Solution {
	root := mode.RootCandidate(sudok)
	b := &backtracker{
		sudok:     sudok,
		solutions: make(chan sudoku.Solution, 1),
	}
	go func() {
		defer close(b.solutions)
		b.backtrack(ctx, root)
	}()
	return b.solutions
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

func (b *backtracker) backtrack(ctx context.Context, candidate Candidate) {
	if ctx.Err() != nil {
		return
	}
	if b.sudok.IsViolated(candidate) {
		return
	}
	if b.sudok.IsSolved(candidate) {
		select {
		case <-ctx.Done():
		case b.solutions <- candidate:
		}
		return
	}
	for _, nextCandidate := range candidate.NextCandidates() {
		b.backtrack(ctx, nextCandidate)
	}
}

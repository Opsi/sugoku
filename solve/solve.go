package solve

import (
	"fmt"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

var (
	errInvalidSolution = fmt.Errorf("invalid solution")
	errNoSolutionFound = fmt.Errorf("no solution found")
)

type solver struct {
	sudok       sudoku.Sudoku
	solution    sudoku.Solution
	fixedValues map[sudoku.Coordinate]struct{}
}

// Solve takes a sudoku and returns a solution, or an error if no solution
// exists for the given sudoku or if the sudoku is invalid.
func Solve(sudok sudoku.Sudoku) (sudoku.Solution, error) {
	// We run a simple backtracking algorithm here.
	// We start with an empty solution, and then we try to fill in the first
	// coordinate with the first possible value. If that works, we try to fill
	// in the second coordinate with the first possible value, and so on.
	// If we reach a point where we can't fill in a coordinate with any of the
	// possible values, we backtrack to the previous coordinate and try the next
	// possible value there.
	// If we reach a point where we can't fill in a coordinate with any of the
	// possible values, and there are no previous coordinates, then we have
	// exhausted all possible solutions and we return an error.
	// If we reach a point where we have filled in all coordinates, then we
	// return the solution.
	s := solver{
		sudok:       sudok,
		solution:    make(sudoku.Solution, len(sudok.Coordinates)),
		fixedValues: make(map[sudoku.Coordinate]struct{}),
	}
	for _, constr := range sudok.Constraints {
		fvc, ok := constr.(constraint.FixedValueConstraint)
		if !ok {
			continue
		}
		s.solution[fvc.Coordinate] = fvc.Value
		s.fixedValues[fvc.Coordinate] = struct{}{}
	}
	err := s.solve(0)
	if err != nil {
		return nil, err
	}
	return s.solution, nil
}

func (s solver) solve(coorIndex int) error {
	if coorIndex >= len(s.sudok.Coordinates) {
		// we have filled in all coordinates and just need to check if the solution is valid
		isSolved := s.sudok.IsSolved(s.solution)
		if !isSolved {
			return errInvalidSolution
		}
		return nil
	}
	coordinate := s.sudok.Coordinates[coorIndex]
	if _, ok := s.fixedValues[coordinate]; ok {
		// this coordinate is fixed, so we don't need to try any values
		return s.solve(coorIndex + 1)
	}
	for _, value := range s.sudok.PossibleValues {
		// try to fill in the coordinate with the value
		s.solution[coordinate] = value
		isViolated := s.sudok.IsViolated(s.solution)
		if isViolated {
			continue
		}
		err := s.solve(coorIndex + 1)
		if err != nil {
			continue
		}
		// we found a solution
		return nil
	}
	// no possible value created a solution
	delete(s.solution, coordinate)
	return errNoSolutionFound
}

type Result uint8

const (
	ResultValidAndSolved = 1 + iota
	ResultValidAndNotSolved
	ResultInvalid
)

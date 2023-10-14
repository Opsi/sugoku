package solve

import (
	"fmt"
	"sudoku-solver/sudoku"
)

var errInvalidSolution = fmt.Errorf("invalid solution")

type solver struct {
	sudok           sudoku.Sudoku
	state           state
	coordinateOrder []sudoku.Coordinate
	deepest         int
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
		sudok:           sudok,
		state:           initState(sudok),
		coordinateOrder: sudok.Coordinates,
		deepest:         0,
	}
	success := s.solve(0)
	if !success {
		return nil, errInvalidSolution
	}
	return s.state, nil
}

func (s *solver) solve(coorIndex int) bool {
	if coorIndex > s.deepest {
		s.deepest = coorIndex
		fmt.Println("deepest:", s.deepest)
	}
	if coorIndex >= len(s.coordinateOrder) {
		// we have filled in all coordinates and just need to check if the solution is valid
		return s.sudok.IsSolved(s.state)
	}
	coordinate := s.coordinateOrder[coorIndex]
	coordState, ok := s.state[coordinate]
	if !ok {
		// this should never happen
		panic(fmt.Sprintf("coordinate %v not found in state", coordinate))
	}
	if coordState.HasValue {
		// this coordinate is fixed, so we don't need to try any values
		return s.solve(coorIndex + 1)
	}
	currentPossibilities := coordState.Possibilities
	for _, value := range currentPossibilities {
		// try to fill in the coordinate with the value
		s.state[coordinate] = coordinateState{
			HasValue:      true,
			Value:         value,
			Possibilities: nil,
		}
		isViolated := s.sudok.IsViolated(s.state)
		if isViolated {
			continue
		}
		success := s.solve(coorIndex + 1)
		if !success {
			continue
		}
		// we found a solution
		return true
	}
	// no possible value created a solution
	s.state[coordinate] = coordinateState{
		HasValue:      false,
		Value:         0,
		Possibilities: currentPossibilities,
	}
	return false
}

type Result uint8

const (
	ResultValidAndSolved = 1 + iota
	ResultValidAndNotSolved
	ResultInvalid
)

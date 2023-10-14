package solve

import (
	"fmt"
	"math"
	"sort"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

var (
	errInvalidSolution = fmt.Errorf("invalid solution")
	errNoSolutionFound = fmt.Errorf("no solution found")
)

type solver struct {
	sudok           sudoku.Sudoku
	solution        sudoku.Solution
	fixedValues     map[sudoku.Coordinate]struct{}
	coordinateOrder []sudoku.Coordinate
	deepest         int
}

// CreateCoordinateOrder creates an order for coordinates to be filled in by the
// backtracking algorithm. The order is based on the number of constraints that
// each coordinate is involved in. Coordinates that are involved in more
// constraints are handled first to fail faster.
func CreateCoordinateOrder(sudok sudoku.Sudoku) []sudoku.Coordinate {
	return sudok.Coordinates
	constraintCount := make(map[sudoku.Coordinate]int, len(sudok.Coordinates))
	for _, coor := range sudok.Coordinates {
		constraintCount[coor] = 0
	}

	// handle fixed value constraints first
	fixedValues := make(map[sudoku.Coordinate]struct{})
	for _, constr := range sudok.Constraints {
		fvc, ok := constr.(constraint.FixedValueConstraint)
		if !ok {
			continue
		}
		// fixed value constraints are already satisfied, so they should
		// be "handled" (skipped) first
		constraintCount[fvc.Coordinate] += math.MaxInt32
		fixedValues[fvc.Coordinate] = struct{}{}
	}

	// handle other constraints
	for _, constr := range sudok.Constraints {
		if _, ok := constr.(constraint.FixedValueConstraint); ok {
			// fixed value constraints were already handled above
			continue
		}
		// for each fixed value this constraint contains, add 1 to the
		// constraint count for each coordinate in the constraint
		fvcCount := 0
		constrainedCoords := constr.ConstrainedCoordinates()
		for _, coor := range constrainedCoords {
			if _, ok := fixedValues[coor]; ok {
				fvcCount++
			}
		}
		for _, coor := range constrainedCoords {
			constraintCount[coor] += fvcCount
		}
	}
	// we want to handle coordinates with the highest constraint count first
	return SortByConstraintCount(constraintCount)
}

func SortByConstraintCount(m map[sudoku.Coordinate]int) []sudoku.Coordinate {
	if len(m) == 0 {
		return nil
	}
	type kv struct {
		Key             sudoku.Coordinate
		ConstraintCount int
	}
	sorted := make([]kv, 0, len(m))
	for k, v := range m {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].ConstraintCount != sorted[j].ConstraintCount {
			return sorted[i].ConstraintCount < sorted[j].ConstraintCount
		}
		if sorted[i].Key.Row != sorted[j].Key.Row {
			return sorted[i].Key.Row < sorted[j].Key.Row
		}
		return sorted[i].Key.Col < sorted[j].Key.Col
	})
	keys := make([]sudoku.Coordinate, 0, len(m))
	for _, kv := range sorted {
		keys = append(keys, kv.Key)
	}
	return keys
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
	coordinateOrder := CreateCoordinateOrder(sudok)
	s := solver{
		sudok:           sudok,
		solution:        make(sudoku.Solution, len(sudok.Coordinates)),
		fixedValues:     make(map[sudoku.Coordinate]struct{}),
		coordinateOrder: coordinateOrder,
		deepest:         0,
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

func (s *solver) solve(coorIndex int) error {
	if coorIndex > s.deepest {
		s.deepest = coorIndex
		fmt.Println("deepest:", s.deepest)
	}
	if coorIndex >= len(s.coordinateOrder) {
		// we have filled in all coordinates and just need to check if the solution is valid
		isSolved := s.sudok.IsSolved(s.solution)
		if !isSolved {
			return errInvalidSolution
		}
		return nil
	}
	coordinate := s.coordinateOrder[coorIndex]
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

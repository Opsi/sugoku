package solve

import (
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

type coordinateState struct {
	HasValue      bool
	Value         int
	Possibilities []int
}

type state map[sudoku.Coordinate]coordinateState

var _ sudoku.Solution = state{}

func (s state) Get(coordinate sudoku.Coordinate) (int, bool) {
	coorState, ok := s[coordinate]
	if !ok {
		return 0, false
	}
	if !coorState.HasValue {
		return 0, false
	}
	return coorState.Value, true
}

func initState(sudok sudoku.Sudoku) state {
	m := make(map[sudoku.Coordinate]coordinateState, len(sudok.Coordinates))
	for _, coor := range sudok.Coordinates {
		m[coor] = coordinateState{
			HasValue:      false,
			Value:         0,
			Possibilities: sudok.PossibleValues,
		}
	}

	// as a simple proof of concept, we will fill in all fixed values
	for _, constr := range sudok.Constraints {
		fvc, ok := constr.(constraint.FixedValueConstraint)
		if !ok {
			continue
		}
		m[fvc.Coordinate] = coordinateState{
			HasValue:      true,
			Value:         fvc.Value,
			Possibilities: nil,
		}
	}

	// now we will remove all values that are not possible for each coordinate
	for _, constr := range sudok.Constraints {
		nrc, ok := constr.(constraint.NoRepeatConstraint)
		if !ok {
			continue
		}
		fixedValues := make(map[int]struct{})
		for _, coor := range nrc.ConstrainedCoordinates() {
			coorState := m[coor]
			if coorState.HasValue {
				fixedValues[coorState.Value] = struct{}{}
			}
		}
		if len(fixedValues) == 0 {
			// no fixed values in this constraint, so we can continue
			continue
		}
		for _, coor := range nrc.ConstrainedCoordinates() {
			coorState := m[coor]
			if !coorState.HasValue {
				newPossibilities := make([]int, 0, len(coorState.Possibilities))
				for _, value := range coorState.Possibilities {
					if _, ok := fixedValues[value]; !ok {
						newPossibilities = append(newPossibilities, value)
					}
				}
				coorState.Possibilities = newPossibilities
				m[coor] = coorState
			}
		}
	}

	return m
}

package backtrack

import (
	"fmt"
	"maps"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

type coordinateState struct {
	HasValue      bool
	Value         int
	Possibilities []int
}

type simpleCandidate struct {
	sudok           sudoku.Sudoku
	coordinateOrder []sudoku.Coordinate

	coordinateIndex int
	state           map[sudoku.Coordinate]coordinateState
}

var _ Candidate = (*simpleCandidate)(nil)

func rootSimple(sudok sudoku.Sudoku) Candidate {
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

	return &simpleCandidate{
		sudok:           sudok,
		coordinateOrder: sudok.Coordinates,
		coordinateIndex: 0,
		state:           m,
	}
}

func (c *simpleCandidate) IsBroken() bool {
	return c.sudok.IsViolated(c)
}

func (c *simpleCandidate) IsSolved() bool {
	return c.sudok.IsSolved(c)
}

func (c *simpleCandidate) NextCandidates() []Candidate {
	// if this candidate belongs to the last coordinate then we are done
	if c.coordinateIndex >= len(c.coordinateOrder)-1 {
		return nil
	}
	coord := c.coordinateOrder[c.coordinateIndex]
	coordState, ok := c.state[coord]
	if !ok {
		// this should never happen
		// TODO: use slog here
		panic(fmt.Sprintf("coordinate %v not found in state", coord))
	}
	if coordState.HasValue {
		// this coordinate is fixed, so we don't need to try any values
		return []Candidate{&simpleCandidate{
			sudok:           c.sudok,
			coordinateOrder: c.coordinateOrder,
			coordinateIndex: c.coordinateIndex + 1,
			state:           c.state,
		}}
	}
	currentPossibilities := coordState.Possibilities
	nextCandidates := make([]Candidate, 0, len(currentPossibilities))
	for _, value := range currentPossibilities {
		// try to fill in the coordinate with the value
		newState := maps.Clone(c.state)
		newState[coord] = coordinateState{
			HasValue:      true,
			Value:         value,
			Possibilities: nil,
		}
		nextCandidates = append(nextCandidates, &simpleCandidate{
			sudok:           c.sudok,
			coordinateOrder: c.coordinateOrder,
			coordinateIndex: c.coordinateIndex + 1,
			state:           newState,
		})
	}
	return nextCandidates
}

func (s simpleCandidate) Get(coordinate sudoku.Coordinate) (int, bool) {
	coorState, ok := s.state[coordinate]
	if !ok {
		return 0, false
	}
	if !coorState.HasValue {
		return 0, false
	}
	return coorState.Value, true
}

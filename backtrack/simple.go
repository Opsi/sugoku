package backtrack

import (
	"fmt"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

type simpleCandidate struct {
	cellsState
	sudok           sudoku.Sudoku
	coordinateOrder []sudoku.Coordinate

	coordinateIndex int
}

var _ Candidate = (*simpleCandidate)(nil)

func rootSimple(sudok sudoku.Sudoku) (Candidate, error) {
	candidate := &simpleCandidate{
		cellsState:      make(cellsState, len(sudok.Coordinates)),
		sudok:           sudok,
		coordinateOrder: sudok.Coordinates,
		coordinateIndex: 0,
	}
	for _, coor := range sudok.Coordinates {
		candidate.cellsState[coor] = VariableCellState(sudok.PossibleValues)
	}

	// as a simple proof of concept, we will fill in all fixed values
	for _, constr := range sudok.Constraints {
		fvc, ok := constr.(constraint.FixedValueConstraint)
		if !ok {
			continue
		}

		candidate.cellsState[fvc.Coordinate] = FixedCellState(fvc.Value)
	}

	// now we will remove all values that are not possible for each coordinate
	for _, constr := range sudok.Constraints {
		nrc, ok := constr.(constraint.NoRepeatConstraint)
		if !ok {
			continue
		}
		fixedValues := make([]int, 0)
		for _, coor := range nrc.ConstrainedCoordinates() {
			coorState := candidate.cellsState[coor]
			if coorState.HasValue {
				fixedValues = append(fixedValues, coorState.Value)
			}
		}
		if len(fixedValues) == 0 {
			// no fixed values in this constraint, so we can continue
			continue
		}
		for _, coor := range nrc.ConstrainedCoordinates() {
			coorState := candidate.cellsState[coor]
			updated, ok := coorState.WithRemovedPossibilities(fixedValues...)
			if !ok {
				// this coordinate is no longer solvable
				return nil, fmt.Errorf("coordinate %v is no longer solvable", coor)
			}
			candidate.cellsState[coor] = updated
		}
	}

	return candidate, nil
}

func (c *simpleCandidate) NextCandidates() []Candidate {
	// if this candidate belongs to the last coordinate then we are done
	if c.coordinateIndex >= len(c.coordinateOrder)-1 {
		return nil
	}
	coord := c.coordinateOrder[c.coordinateIndex]
	cell, ok := c.cellsState[coord]
	if !ok {
		// this should never happen
		// TODO: use slog here
		panic(fmt.Sprintf("coordinate %v not found in state", coord))
	}
	if cell.HasValue {
		// this coordinate is fixed, so we don't need to try any values
		return []Candidate{&simpleCandidate{
			cellsState:      c.cellsState.Copy(),
			sudok:           c.sudok,
			coordinateOrder: c.coordinateOrder,
			coordinateIndex: c.coordinateIndex + 1,
		}}
	}
	currentPossibilities := cell.Possibilities
	nextCandidates := make([]Candidate, 0, len(currentPossibilities))
	for _, value := range currentPossibilities {
		// try to fill in the coordinate with the value
		newState := c.cellsState.Copy()
		newState[coord] = cellState{
			HasValue:      true,
			Value:         value,
			Possibilities: nil,
		}
		nextCandidates = append(nextCandidates, &simpleCandidate{
			cellsState:      newState,
			sudok:           c.sudok,
			coordinateOrder: c.coordinateOrder,
			coordinateIndex: c.coordinateIndex + 1,
		})
	}
	return nextCandidates
}

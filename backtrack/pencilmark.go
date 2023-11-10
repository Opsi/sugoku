package backtrack

import (
	"fmt"
	"slices"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"
)

type pencilmarkCandidate struct {
	sudok sudoku.Sudoku

	coordinateIndex int
	state           map[sudoku.Coordinate]cellState
}

var _ Candidate = (*pencilmarkCandidate)(nil)

func (c *pencilmarkCandidate) FillIn(coordinate sudoku.Coordinate, value int) error {
	if _, ok := c.state[coordinate]; !ok {
		return fmt.Errorf("coordinate %v not found in state", coordinate)
	}
	c.state[coordinate] = cellState{
		HasValue:      true,
		Value:         value,
		Possibilities: nil,
	}
	for _, constr := range c.sudok.Constraints {
		nrc, ok := constr.(constraint.NoRepeatConstraint)
		if !ok {
			continue
		}
		// we need to check if this constraint contains the coordinate
		if !slices.Contains(nrc.ConstrainedCoordinates(), coordinate) {
			continue
		}
		for _, coor := range nrc.ConstrainedCoordinates() {
			coorState := c.state[coor]
			stillSolvable := coorState.RemovePossibilities(value)
			if !stillSolvable {
				return fmt.Errorf("coordinate %v is no longer solvable", coor)
			}
			c.state[coor] = coorState
		}
	}
	return nil
}

func rootPencilMark(sudok sudoku.Sudoku) (Candidate, error) {
	m := make(map[sudoku.Coordinate]cellState, len(sudok.Coordinates))
	for _, coor := range sudok.Coordinates {
		m[coor] = cellState{
			HasValue:      false,
			Value:         0,
			Possibilities: sudok.PossibleValues,
		}
	}

	candidate := &pencilmarkCandidate{
		sudok:           sudok,
		coordinateIndex: 0,
		state:           m,
	}

	// first we fill in all fixed values
	for _, constr := range sudok.Constraints {
		fvc, ok := constr.(constraint.FixedValueConstraint)
		if !ok {
			continue
		}

		err := candidate.FillIn(fvc.Coordinate, fvc.Value)
		if err != nil {
			return nil, fmt.Errorf("fill in fixed value %d at %v: %w", fvc.Value, fvc.Coordinate, err)
		}
	}

	return candidate, nil
}

func (c *pencilmarkCandidate) IsBroken() bool {
	return c.sudok.IsViolated(c)
}

func (c *pencilmarkCandidate) IsSolved() bool {
	return c.sudok.IsSolved(c)
}

func (c *pencilmarkCandidate) NextCandidates() []Candidate {
	// if this candidate belongs to the last coordinate then we are done
	if c.coordinateIndex >= len(c.sudok.Coordinates)-1 {
		return nil
	}
	coord := c.sudok.Coordinates[c.coordinateIndex]
	cell, ok := c.state[coord]
	if !ok {
		// this should never happen
		// TODO: use slog here
		panic(fmt.Sprintf("coordinate %v not found in state", coord))
	}
	if cell.HasValue {
		// this coordinate is fixed, so we don't need to try any values
		return []Candidate{&pencilmarkCandidate{
			sudok:           c.sudok,
			coordinateIndex: c.coordinateIndex + 1,
			state:           c.state,
		}}
	}

	currentPossibilities := cell.Possibilities
	nextCandidates := make([]Candidate, 0, len(currentPossibilities))
	for _, value := range currentPossibilities {
		// try to fill in the coordinate with the value
		newState := make(map[sudoku.Coordinate]cellState, len(c.state))
		for coor, state := range c.state {
			newState[coor] = state.Copy()
		}

		newCandidate := &pencilmarkCandidate{
			sudok:           c.sudok,
			coordinateIndex: c.coordinateIndex + 1,
			state:           newState,
		}
		err := newCandidate.FillIn(coord, value)
		if err != nil {
			// filling in this value makes the sudoku unsolvable
			continue
		}
		nextCandidates = append(nextCandidates, newCandidate)
	}
	return nextCandidates
}

func (s pencilmarkCandidate) Get(coordinate sudoku.Coordinate) (int, bool) {
	coorState, ok := s.state[coordinate]
	if !ok {
		return 0, false
	}
	if !coorState.HasValue {
		return 0, false
	}
	return coorState.Value, true
}

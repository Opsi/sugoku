package backtrack

import (
	"fmt"
	"slices"
	"sudoku-solver/constraint"
	"sudoku-solver/sudoku"

	"golang.org/x/exp/maps"
)

type pencilmarkCandidate struct {
	cellsState
	sudok sudoku.Sudoku

	coordinateIndex int
}

var _ Candidate = (*pencilmarkCandidate)(nil)

func (c *pencilmarkCandidate) FillIn(coordinate sudoku.Coordinate, value int) error {
	if _, ok := c.cellsState[coordinate]; !ok {
		return fmt.Errorf("coordinate %v not found in state", coordinate)
	}
	c.cellsState[coordinate] = FixedCellState(value)
	for _, constr := range c.sudok.Constraints {
		switch constr := constr.(type) {
		case constraint.NoRepeatConstraint:
			if err := c.updateWithNoRepeatConstraint(constr, coordinate, value); err != nil {
				return fmt.Errorf("update with no repeat constraint: %w", err)
			}
		case constraint.SameSumConstraint:
			if err := c.updateWithSameSumConstraint(constr, coordinate, value); err != nil {
				return fmt.Errorf("update with same sum constraint: %w", err)
			}
		default:
		}
	}
	return nil
}

func (c *pencilmarkCandidate) updateWithNoRepeatConstraint(constr constraint.NoRepeatConstraint, coordinate sudoku.Coordinate, value int) error {
	// In a NoRepeatConstraint, we need to remove the value from the possibilities
	// of all other coordinates in the constraint.
	if !slices.Contains(constr.ConstrainedCoordinates(), coordinate) {
		return nil
	}
	for _, coor := range constr.ConstrainedCoordinates() {
		coorState := c.cellsState[coor]
		updated, stillSolvable := coorState.WithRemovedPossibilities(value)
		if !stillSolvable {
			return fmt.Errorf("coordinate %v is no longer solvable", coor)
		}
		c.cellsState[coor] = updated
	}
	return nil
}

func (c *pencilmarkCandidate) updateWithSameSumConstraint(constr constraint.SameSumConstraint, coordinate sudoku.Coordinate, value int) error {
	// Since this is easily expandable we just start with the case
	// where we have an arrow and something was written on any of the
	// path coordinates.
	if len(constr.Coordinates1) != 1 {
		// this is not an arrow constraint
		return nil
	}
	circle := constr.Coordinates1[0]
	if circle == coordinate {
		// we filled in the circle coordinate, so for now we don't
		// care about this case
		return nil
	}
	circleState, ok := c.cellsState[circle]
	if !ok {
		// this should never happen
		panic(fmt.Sprintf("coordinate %v not found in state", circle))
	}
	if circleState.HasValue {
		// we filled in the circle coordinate already, so for now we don't
		// care about this case
		return nil
	}
	path := constr.Coordinates2
	if !slices.Contains(path, coordinate) {
		// the coordinate we filled in is not on the path, so for now
		// we don't care about this case
		return nil
	}
	// we filled in a coordinate on the path, so we need can maybe remove
	// some possibilities from the circle coordinate
	// so we just calculate all the possible sums and remove all values
	// that are not possible anymore
	possibilitiesSlice := Map(path, func(coor sudoku.Coordinate) []int {
		coorState := c.cellsState[coor]
		if coorState.HasValue {
			return []int{coorState.Value}
		}
		return coorState.Possibilities
	})
	possibleSums := make(map[int]struct{})
	for _, valueCombination := range Combinations(possibilitiesSlice) {
		sum := Reduce(valueCombination, 0, func(acc, value int) int {
			return acc + value
		})
		possibleSums[sum] = struct{}{}
	}
	newCircleState, stillSolvable := c.cellsState[circle].WithConstrainedPossibilities(maps.Keys(possibleSums)...)
	if !stillSolvable {
		return fmt.Errorf("arrow circle coordinate %v is no longer solvable", circle)
	}
	c.cellsState[circle] = newCircleState
	return nil
}

func rootPencilMark(sudok sudoku.Sudoku) (Candidate, error) {
	candidate := &pencilmarkCandidate{
		cellsState:      make(cellsState, len(sudok.Coordinates)),
		sudok:           sudok,
		coordinateIndex: 0,
	}
	for _, coor := range sudok.Coordinates {
		candidate.cellsState[coor] = VariableCellState(sudok.PossibleValues)
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

func (c *pencilmarkCandidate) NextCandidates() []Candidate {
	// if this candidate belongs to the last coordinate then we are done
	if c.coordinateIndex >= len(c.sudok.Coordinates)-1 {
		return nil
	}
	coord := c.sudok.Coordinates[c.coordinateIndex]
	cell, ok := c.cellsState[coord]
	if !ok {
		// this should never happen
		// TODO: use slog here
		panic(fmt.Sprintf("coordinate %v not found in state", coord))
	}
	if cell.HasValue {
		// this coordinate is fixed, so we don't need to try any values
		return []Candidate{&pencilmarkCandidate{
			cellsState:      c.cellsState.Copy(),
			sudok:           c.sudok,
			coordinateIndex: c.coordinateIndex + 1,
		}}
	}

	currentPossibilities := cell.Possibilities
	nextCandidates := make([]Candidate, 0, len(currentPossibilities))
	for _, value := range currentPossibilities {
		// try to fill in the coordinate with the value
		newCandidate := &pencilmarkCandidate{
			cellsState:      c.cellsState.Copy(),
			sudok:           c.sudok,
			coordinateIndex: c.coordinateIndex + 1,
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

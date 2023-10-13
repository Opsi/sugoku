package sudoku

import "fmt"

type Coordinate struct {
	Row, Col int
}

type ConstraintResult uint8

const (
	ConstraintResultValidAndSolved = 1 + iota
	ConstraintResultValidAndNotSolved
	ConstraintResultInvalid
)

type Constraint interface {
	Check(Solution) ConstraintResult
}

type Sudoku struct {
	Coordinates    []Coordinate
	PossibleValues []int
	Constraints    []Constraint
}

type Solution = map[Coordinate]int

func (s Sudoku) Solve() (Solution, error) {
	return Solution{}, fmt.Errorf("not implemented")
}

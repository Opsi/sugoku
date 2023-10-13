package sudoku

import "fmt"

type Coordinate struct {
	Row, Col int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("%d-%d", c.Row, c.Col)
}

type Constraint interface {
	IsViolated(Solution) bool
}

type Sudoku struct {
	Coordinates    []Coordinate
	PossibleValues []int
	Constraints    []Constraint
}

var _ Constraint = Sudoku{}

type Solution = map[Coordinate]int

func (s Sudoku) IsViolated(solution Solution) bool {
	for _, constraint := range s.Constraints {
		if constraint.IsViolated(solution) {
			return true
		}
	}
	return false
}

func (s Sudoku) IsSolved(solution Solution) bool {
	if s.IsViolated(solution) {
		return false
	}
	for _, coordinate := range s.Coordinates {
		if _, ok := solution[coordinate]; !ok {
			return false
		}
	}
	return true
}

package sudoku

import "fmt"

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

func (s Sudoku) Check(solution Solution) error {
	for _, constraint := range s.Constraints {
		if constraint.IsViolated(solution) {
			return fmt.Errorf("solution violates constraint %v", constraint)
		}
	}
	for _, coordinate := range s.Coordinates {
		if _, ok := solution[coordinate]; !ok {
			return fmt.Errorf("solution is missing coordinate %v", coordinate)
		}
	}
	return nil
}

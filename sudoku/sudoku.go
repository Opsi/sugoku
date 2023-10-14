package sudoku

import "fmt"

type Constraint interface {
	IsViolated(Solution) bool
	ConstrainedCoordinates() []Coordinate
}

type Sudoku struct {
	Coordinates    []Coordinate
	PossibleValues []int
	Constraints    []Constraint
}

type Solution interface {
	Get(Coordinate) (int, bool)
}

type solution map[Coordinate]int

func (s solution) Get(coordinate Coordinate) (int, bool) {
	value, ok := s[coordinate]
	return value, ok
}

func MapSolution(m map[Coordinate]int) Solution {
	return solution(m)
}

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
		if _, ok := solution.Get(coordinate); !ok {
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
		if _, ok := solution.Get(coordinate); !ok {
			return fmt.Errorf("solution is missing coordinate %v", coordinate)
		}
	}
	return nil
}

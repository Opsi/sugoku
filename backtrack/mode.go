package backtrack

import (
	"fmt"
	"strings"
	"sudoku-solver/sudoku"
)

type Mode string

const (
	ModeSimple Mode = "simple"
)

func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(s) {
	case string(ModeSimple):
		return ModeSimple, nil
	default:
		return "", fmt.Errorf("unknown mode: %s", s)
	}
}

func MustParseMode(s string) Mode {
	mode, err := ParseMode(s)
	if err != nil {
		panic(err)
	}
	return mode
}

func (m Mode) String() string {
	return string(m)
}

func (m Mode) RootCandidate(sudok sudoku.Sudoku) Candidate {
	switch m {
	case ModeSimple:
		return rootSimple(sudok)
	default:
		panic(fmt.Errorf("unknown mode: %s", m))
	}
}

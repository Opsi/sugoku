package backtrack

import (
	"fmt"
	"strings"
	"sudoku-solver/sudoku"
)

type Mode string

const (
	ModeSimple     Mode = "simple"
	ModePencilMark Mode = "pencilmark"
)

func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(s) {
	case string(ModeSimple):
		return ModeSimple, nil
	case string(ModePencilMark):
		return ModePencilMark, nil
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
	case ModePencilMark:
		root, err := rootPencilMark(sudok)
		if err != nil {
			// TODO: better error handling
			panic(err)
		}
		return root
	default:
		panic(fmt.Errorf("unknown mode: %s", m))
	}
}

package main

import (
	"fmt"
	"os"
	"strings"
	"sudoku-solver/solve"
	"sudoku-solver/sudoku"
	"sudoku-solver/sudokuio"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// first argument is the input file
	if len(os.Args) < 2 {
		return fmt.Errorf("no input file specified")
	}
	inputFilePath := os.Args[1]

	// read the input file
	inputBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("read input file: %v", err)
	}

	var sudok *sudoku.Sudoku
	if strings.HasSuffix(inputFilePath, ".txt") {
		// parse the input file into a string
		sudok, err = sudokuio.ParseString(string(inputBytes))
		if err != nil {
			return fmt.Errorf("parse txt file: %v", err)
		}
	} else if strings.HasSuffix(inputFilePath, ".json") {
		// parse the input file as JSON
		sudok, err = sudokuio.ParseJSON(inputBytes)
		if err != nil {
			return fmt.Errorf("parse json file: %v", err)
		}
	} else {
		return fmt.Errorf("unknown input file type")
	}

	// solve the sudoku
	solution, err := solve.Solve(*sudok)
	if err != nil {
		return fmt.Errorf("solve sudoku: %v", err)
	}

	// print the solution
	lastRow := -9999
	for _, coord := range sudok.Coordinates {
		value, ok := solution.Get(coord)
		if !ok {
			return fmt.Errorf("solution missing value for coordinate %v", coord)
		}
		if lastRow != coord.Row {
			fmt.Println()
			lastRow = coord.Row
		}
		fmt.Printf("%d", value)
	}
	return nil
}

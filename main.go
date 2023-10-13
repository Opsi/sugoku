package main

import (
	"fmt"
	"os"
	"sudoku-solver/solve"
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

	// parse the input file into a string
	sudok, err := sudokuio.ParseString(string(inputBytes))
	if err != nil {
		return fmt.Errorf("parse input file: %v", err)
	}

	// solve the sudoku
	solution, err := solve.Solve(*sudok)
	if err != nil {
		return fmt.Errorf("solve sudoku: %v", err)
	}

	// print the solution
	lastRow := -9999
	for _, coor := range sudok.Coordinates {
		value := solution[coor]
		if lastRow != coor.Row {
			fmt.Println()
			lastRow = coor.Row
		}
		fmt.Printf("%d", value)
	}
	return nil
}

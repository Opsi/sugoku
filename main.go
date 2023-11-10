package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"os/signal"
	"strings"
	"sudoku-solver/backtrack"
	"sudoku-solver/sudoku"
	"sudoku-solver/sudokuio"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	modeString := flag.String("mode", "pencilmark", "mode to use for solving")
	flag.Parse()

	args := flag.Args()
	// first argument is the input file
	if len(args) < 1 {
		return fmt.Errorf("no input file specified")
	}
	inputFilePath := args[0]

	mode, err := backtrack.ParseMode(*modeString)
	if err != nil {
		return fmt.Errorf("parse mode: %w", err)
	}

	// read the input file
	inputBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	var sudok *sudoku.Sudoku
	if strings.HasSuffix(inputFilePath, ".txt") {
		// parse the input file into a string
		sudok, err = sudokuio.ParseString(string(inputBytes))
		if err != nil {
			return fmt.Errorf("parse txt file: %w", err)
		}
	} else if strings.HasSuffix(inputFilePath, ".json") {
		// parse the input file as JSON
		sudok, err = sudokuio.ParseJSON(inputBytes)
		if err != nil {
			return fmt.Errorf("parse json file: %w", err)
		}
	} else {
		return fmt.Errorf("unknown input file type")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// solve the sudoku
	slog.Info("starting to solve", slog.String("mode", mode.String()))
	start := time.Now()
	solutions := backtrack.FindSolutions(ctx, mode, *sudok)
	ticker := time.NewTicker(30 * time.Second)
	solutionCount := 0

	for {
		select {
		case <-ctx.Done():
			slog.Info("stopped solving",
				slog.Int("foundSolutions", solutionCount),
				slog.Duration("searchDuration", time.Since(start)),
			)
			return nil
		case <-ticker.C:
			slog.Info("still searching...",
				slog.Int("foundSolutions", solutionCount),
				slog.Duration("searchDuration", time.Since(start)),
			)
		case solution, more := <-solutions:
			if !more {
				slog.Info("no more solutions",
					slog.Int("foundSolutions", solutionCount),
					slog.Duration("searchDuration", time.Since(start)),
				)
				return nil
			}
			solutionCount++
			slog.Info("found solution",
				slog.Int("foundSolutions", solutionCount),
				slog.Duration("searchDuration", time.Since(start)),
			)

			// print the solution
			lastRow := math.MinInt
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
		}
	}
}

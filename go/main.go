package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <path-to-csv>")
		os.Exit(1)
	}

	csvPath := os.Args[1]
	fmt.Printf("Running benchmark using data from: %s\n", csvPath)
	BenchmarkAggregations(csvPath)
}

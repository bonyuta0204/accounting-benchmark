package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define command-line flag for CSV path
	csvPath := flag.String("csv", "../data/sample_1m.csv", "Path to CSV file")
	flag.Parse()

	fmt.Printf("Running benchmark using data from: %s\n", *csvPath)
	BenchmarkAggregations(*csvPath)
}

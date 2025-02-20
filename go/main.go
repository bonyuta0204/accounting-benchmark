package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define command-line flags
	mode := flag.String("mode", "benchmark", "Mode: generate, process, benchmark")
	csvPath := flag.String("csv", "../data/sample_1m.csv", "Path to CSV file")
	rows := flag.Int("rows", 1000000, "Number of rows for data generation")
	flag.Parse()

	switch *mode {
	case "generate":
		err := GenerateCSV(*csvPath, *rows)
		if err != nil {
			fmt.Println("Error generating CSV:", err)
		} else {
			fmt.Println("CSV generated at", *csvPath)
		}
	case "process":
		// Process aggregations without benchmarking
		ProcessAggregations(*csvPath)
	case "benchmark":
		BenchmarkAggregations(*csvPath)
	default:
		fmt.Println("Invalid mode. Use generate, process, or benchmark.")
	}
}

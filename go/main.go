package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <path-to-csv>")
		os.Exit(1)
	}

	//
	csvPath := os.Args[1]
	var dataFrameRunner = NewDataFrameRunner()
	var naiveRunner = NewNaiveAggregator()
	var aggregators = []Aggregator{dataFrameRunner, naiveRunner}

	for _, aggregator := range aggregators {
		fmt.Println("=== Running benchmark using", aggregator.Name())

		// Load CSV

		start := time.Now()
		if err := aggregator.LoadCSV(csvPath); err != nil {
			fmt.Printf("Error loading CSV: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loading CSV took %s\n", time.Since(start))

		// Month x Date aggregation
		start = time.Now()
		if err := aggregator.Aggregate("Date", "Month"); err != nil {
			fmt.Printf("Error running benchmark: %v\n", err)
		}
		fmt.Printf("Aggregation took %s\n", time.Since(start))

		// Write to CSV
		start = time.Now()
		outPath := fmt.Sprintf("%s_%s.csv", csvPath, aggregator.Name())
		if err := aggregator.WriteToCSV(outPath); err != nil {
			fmt.Printf("Error writing to CSV: %v\n", err)
		}
		fmt.Printf("Writing to CSV took %s\n", time.Since(start))
	}
}

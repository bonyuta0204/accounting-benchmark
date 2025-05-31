package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <path-to-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// Detect file format
	var isParquet bool
	switch ext {
	case ".csv":
		isParquet = false
	case ".parquet":
		isParquet = true
	default:
		fmt.Printf("Unsupported file format: %s. Only .csv and .parquet files are supported.\n", ext)
		os.Exit(1)
	}
	
	fmt.Printf("Processing %s file: %s\n", ext[1:], filePath)
	
	var dataFrameRunner = NewDataFrameRunner()
	var naiveRunner = NewNaiveAggregator()
	var duckdbRunner = NewDuckDBAggregator()
	var aggregators = []Aggregator{dataFrameRunner, naiveRunner, duckdbRunner}

	for _, aggregator := range aggregators {
		fmt.Printf("\n=== Running benchmark using %s ===\n", aggregator.Name())

		// Clean up DuckDB connection after use
		if duckdb, ok := aggregator.(*DuckDBAggregator); ok {
			defer duckdb.Close()
		}

		// Load data
		start := time.Now()
		var loadErr error
		if isParquet {
			loadErr = aggregator.LoadParquet(filePath)
		} else {
			loadErr = aggregator.LoadCSV(filePath)
		}
		
		if loadErr != nil {
			fmt.Printf("Error loading %s: %v\n", ext[1:], loadErr)
			continue
		}
		loadTime := time.Since(start)
		fmt.Printf("Loading %s took %s\n", ext[1:], loadTime)

		// Month x Date aggregation
		start = time.Now()
		if err := aggregator.Aggregate("Date", "Month"); err != nil {
			fmt.Printf("Error running aggregation: %v\n", err)
			continue
		}
		aggregateTime := time.Since(start)
		fmt.Printf("Aggregation took %s\n", aggregateTime)

		// Write to CSV
		start = time.Now()
		outPath := fmt.Sprintf("%s_%s.csv", strings.TrimSuffix(filePath, ext), aggregator.Name())
		if err := aggregator.WriteToCSV(outPath); err != nil {
			fmt.Printf("Error writing to CSV: %v\n", err)
			continue
		}
		writeTime := time.Since(start)
		fmt.Printf("Writing to CSV took %s\n", writeTime)
		
		// Print summary
		fmt.Printf("\nSummary for %s:\n", aggregator.Name())
		fmt.Printf("  Load time:      %s\n", loadTime)
		fmt.Printf("  Aggregate time: %s\n", aggregateTime)
		fmt.Printf("  Write time:     %s\n", writeTime)
		fmt.Printf("  Total time:     %s\n", loadTime + aggregateTime + writeTime)
	}
}

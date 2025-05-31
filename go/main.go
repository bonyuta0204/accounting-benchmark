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
	var dataFrameRunner = NewDataFrameRunner(csvPath)
	var naiveRunner = NewNaiveRunner(csvPath)
	var runners = []AggregationRunner{dataFrameRunner, naiveRunner}

	for _, runner := range runners {
		fmt.Println("=== Running benchmark using", runner.Name())
		start := time.Now()
		if err := runner.Run(); err != nil {
			fmt.Printf("Error running benchmark: %v\n", err)
		}
		fmt.Printf("Total time: %s\n", time.Since(start))
	}
}

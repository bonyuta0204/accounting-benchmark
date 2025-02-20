package main

import (
    "fmt"
    "time"
)

func BenchmarkAggregations(csvPath string) {
    start := time.Now()
    ProcessAggregations(csvPath)
    elapsed := time.Since(start)
    fmt.Printf("Total processing time: %s\n", elapsed)
}


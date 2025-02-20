package main

import (
	"fmt"
	"time"
)

func BenchmarkAggregations(csvPath string) {
	// Benchmark Account × Monthly Aggregation
	start := time.Now()
	ok := ProcessAccountMonthAggregation(csvPath)
	elapsed := time.Since(start)
	fmt.Printf("Account × Monthly Aggregation took: %s\n", elapsed)
	if !ok {
		fmt.Println("Error in Account × Monthly Aggregation")
	}

	// Benchmark Department × Monthly Aggregation
	start = time.Now()
	ok = ProcessDepartmentMonthAggregation(csvPath)
	elapsed = time.Since(start)
	fmt.Printf("Department × Monthly Aggregation took: %s\n", elapsed)
	if !ok {
		fmt.Println("Error in Department × Monthly Aggregation")
	}

	// Benchmark Account × Department × Monthly Aggregation
	start = time.Now()
	ok = ProcessAccountDepartmentMonthAggregation(csvPath)
	elapsed = time.Since(start)
	fmt.Printf("Account × Department × Monthly Aggregation took: %s\n", elapsed)
	if !ok {
		fmt.Println("Error in Account × Department × Monthly Aggregation")
	}

	// Create pivot table from account-month aggregation
	ok = CreatePivotTable(csvPath)
	if !ok {
		fmt.Println("Error creating pivot table")
	}
}

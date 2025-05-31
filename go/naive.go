package main

import (
	"fmt"
	"time"
)

type NaiveRunner struct {
	csvPath string
}

func NewNaiveRunner(csvPath string) *NaiveRunner {
	return &NaiveRunner{csvPath: csvPath}
}

func (r *NaiveRunner) Run() error {
	// Load CSV
	start := time.Now()
	fmt.Printf("CSV loading took: %s\n", time.Since(start))
	return nil
}

func (r *NaiveRunner) Name() string {
	return "Naive"
}

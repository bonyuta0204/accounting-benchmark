# Go Implementation for Accounting Data Aggregation Benchmark

## Overview
This Go solution uses the Gota DataFrame library to:
- Generate a sample CSV dataset.
- Perform aggregations:
  - Account × Monthly Aggregation
  - Department × Monthly Aggregation
  - Account × Department × Monthly Aggregation
- Optionally create a pivot table from the results.
- Benchmark the processing time.

## Setup and Execution
1. Ensure you have Go (LTS version) installed.
2. Navigate to the `go` directory.
3. Run `go mod tidy` to download the dependencies.
4. To generate data, run:


go run main.go -mode=generate -csv=../data/sample_1m.csv -rows=1000000

5. To process aggregations (without benchmarking), run:
go run main.go -mode=process -csv=../data/sample_1m.csv

6. To run benchmarks, run:
go run main.go -mode=benchmark -csv=../data/sample_1m.csv


# Go Implementation for Accounting Data Aggregation Benchmark

## Overview

This Go solution uses the [Gota DataFrame](https://github.com/go-gota/gota) library to:
- Perform aggregations:
  - Account × Monthly Aggregation
  - Department × Monthly Aggregation
  - Account × Department × Monthly Aggregation
- Optionally create a pivot table from the results
- Benchmark the processing time

## Setup and Execution

### Prerequisites
- Go (LTS version)
- Sample data (generate using the top-level `generator` package)

### Installation
1. Navigate to the `go` directory
2. Run `go mod tidy` to download the dependencies

### Usage

Run the benchmark:
```bash
# Run with a specific data file
go run . /path/to/data.csv

# Example
go run . ../data/transactions_1m.csv
```

### Command Line Arguments
The program takes a single positional argument:
- Path to the CSV file (required)

### Using Makefile
You can also use the Makefile in the root directory to run benchmarks:
```bash
# Run 1M records benchmark
make benchmark-1m

# Run 10M records benchmark
make benchmark-10m

# Run 100M records benchmark
make benchmark-100m
```

### Output
Results will be saved in the `../results/` directory with the following files:
- `go_account_month.csv`
- `go_department_month.csv`
- `go_account_dept_month.csv`
- `go_pivot_aggregation.csv` (if pivot option is enabled)

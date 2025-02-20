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
# Use default data path (../data/sample_1m.csv)
go run .

# Or specify a custom data path
go run . -csv=../data/sample_10m.csv
```

### Command Line Arguments
- `-csv`: Path to the CSV file (default: "../data/sample_1m.csv")

### Output
Results will be saved in the `../results/` directory with the following files:
- `go_account_month.csv`
- `go_department_month.csv`
- `go_account_dept_month.csv`
- `go_pivot_aggregation.csv` (if pivot option is enabled)

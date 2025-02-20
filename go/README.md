# Go Implementation for Accounting Data Aggregation Benchmark

## Overview

This Go solution uses the [Gota DataFrame](https://github.com/go-gota/gota) library to:
- Generate a sample CSV dataset
- Perform aggregations:
  - Account × Monthly Aggregation
  - Department × Monthly Aggregation
  - Account × Department × Monthly Aggregation
- Optionally create a pivot table from the results
- Benchmark the processing time

## Setup and Execution

### Prerequisites
- Go (LTS version)

### Installation
1. Navigate to the `go` directory
2. Run `go mod tidy` to download the dependencies

### Usage

#### Generate Sample Data
```bash
go run . -mode=generate -csv=../data/sample_1m.csv -rows=1000000
```

#### Process Aggregations
```bash
go run . -mode=process -csv=../data/sample_1m.csv
```

#### Run Benchmarks
```bash
go run . -mode=benchmark -csv=../data/sample_1m.csv
```

### Command Line Arguments
- `-mode`: Operation mode (`generate`, `process`, or `benchmark`)
- `-csv`: Path to the CSV file
- `-rows`: Number of rows to generate (only used with `-mode=generate`)

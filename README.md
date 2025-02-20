# Accounting Data Aggregation Benchmark

## Project Overview
This project compares the performance of Rust's Polars and Go's Gota DataFrame libraries in aggregating accounting data (CSV format with over 1 million rows). The aggregations include:
1. **Account × Monthly Aggregation**
2. **Department × Monthly Aggregation**
3. **Account × Department × Monthly Aggregation**

Results are exported as CSV files (including a pivot table format).

## Directory Structure

```
accounting-benchmark/
├── go/
│   ├── go.mod
│   ├── main.go
│   ├── data_generator.go
│   ├── process.go
│   ├── benchmark.go
│   └── README.md
├── rust/
│   ├── Cargo.toml
│   ├── src/
│   │   ├── main.rs
│   │   ├── data_generator.rs
│   │   ├── process.rs
│   │   └── benchmark.rs
│   └── README.md
├── data/
│   ├── sample_1m.csv
│   ├── sample_10m.csv
│   └── sample_100m.csv
├── results/
│   ├── go_account_month.csv
│   ├── go_department_month.csv
│   ├── go_account_dept_month.csv
│   ├── go_pivot_aggregation.csv
│   ├── rust_account_month.csv
│   ├── rust_department_month.csv
│   ├── rust_account_dept_month.csv
│   └── rust_pivot_aggregation.csv
├── README.md
└── .gitignore
```

## Implementation Flow

1. **Data Preparation**
   - Generate CSV data using the data generator in either Go or Rust.

2. **Aggregation Processing**
   - Each implementation reads the CSV
   - Extracts the month from the Date
   - Performs the required aggregations
   - Optionally pivots the result

3. **Benchmarking**
   - The implementations measure and print the execution time for each aggregation.

4. **Output**
   - Aggregated results and benchmark data are saved as CSV files in the `results/` directory.

## Execution Instructions

### Rust
1. Navigate to the `rust` directory
2. Build with `cargo build` and run with `cargo run`

### Go
1. Navigate to the `go` directory
2. Run `go mod tidy` to fetch dependencies
3. Use:
   ```bash
   go run main.go -mode=<generate|process|benchmark>
   ```
   to generate data, process aggregations, or run benchmarks.

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
├── generator/           # Data generation package (Rust)
│   ├── Cargo.toml
│   └── src/
│       └── main.rs
├── go/                 # Go implementation
│   ├── go.mod
│   ├── main.go
│   ├── process.go
│   ├── benchmark.go
│   └── README.md
├── rust/               # Rust implementation
│   ├── Cargo.toml
│   ├── src/
│   │   ├── main.rs
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
   - Generate CSV data using the dedicated generator package:
     ```bash
     # Generate 1 million rows of data
     cd generator
     cargo run -- -o ../data/sample_1m.csv -r 1000000
     ```

2. **Aggregation Processing**
   - Each implementation reads the CSV
   - Extracts the month from the Date
   - Performs the required aggregations
   - Optionally pivots the result

3. **Benchmarking**
   - The implementations measure and print the execution time for each aggregation.

4. **Output**
   - Aggregated results and benchmark data are saved as CSV files in the `results/` directory.

## Data Generation

The project includes a dedicated Rust-based data generator with the following features:
- Configurable number of rows
- Customizable date range
- Consistent data format for benchmarking
- CLI interface with the following options:
  - `-o, --output`: Output file path
  - `-r, --rows`: Number of rows to generate (default: 1000)
  - `-s, --start-date`: Start date (default: 2020-01-01)
  - `-e, --end-date`: End date (default: 2020-12-31)

Example usage:
```bash
cd generator
# Generate 1 million rows for year 2020
cargo run -- -o ../data/sample_1m.csv -r 1000000

# Generate 10 million rows for year 2023
cargo run -- -o ../data/sample_10m.csv -r 10000000 -s 2023-01-01 -e 2023-12-31
```

## Benchmark Results

Below are the benchmark results comparing Go and Rust implementations:

### Go Implementation
```
Account × Monthly Aggregation took: 1.696857833s
Department × Monthly Aggregation took: 1.683926916s
Account × Department × Monthly Aggregation took: 1.852361s
```

### Rust Implementation (Unoptimized)
```
Account × Monthly Aggregation took: 1.727755958s
Department × Monthly Aggregation took: 1.701820167s
Account × Department × Monthly Aggregation took: 1.794037666s
```

### Rust Implementation (Optimized)
```
Account × Monthly Aggregation took: 183.729083ms
Department × Monthly Aggregation took: 118.934791ms
Account × Department × Monthly Aggregation took: 126.28775ms
```

### Summary
The optimized Rust implementation shows significant performance improvements:
- Account × Monthly aggregation is ~9.2x faster
- Department × Monthly aggregation is ~14.2x faster
- Account × Department × Monthly aggregation is ~14.7x faster

compared to both the Go implementation and unoptimized Rust version.

## Execution Instructions

### Data Generation
1. Navigate to the `generator` directory
2. Run `cargo build` to build the generator
3. Generate data using the CLI options described above

### Rust Implementation
1. Navigate to the `rust` directory
2. Build with `cargo build` and run with `cargo run`

### Go Implementation
1. Navigate to the `go` directory
2. Build with `go build` and run with `go run main.go`

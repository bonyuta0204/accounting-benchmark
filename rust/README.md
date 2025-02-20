# Rust Implementation for Accounting Data Aggregation Benchmark

## Overview

This Rust solution uses the [Polars](https://github.com/pola-rs/polars) library to:
- Perform aggregations:
  - Account × Monthly Aggregation
  - Department × Monthly Aggregation
  - Account × Department × Monthly Aggregation
- Optionally create a pivot table from the results
- Benchmark the processing time

## Setup and Execution

### Prerequisites
- Rust (LTS version)
- Cargo (comes with Rust)
- Sample data (generate using the top-level `generator` package)

### Installation
1. Navigate to the `rust` directory
2. Build the project:
   ```bash
   cargo build --release
   ```

### Usage

Run the benchmark with a specific data file:
```bash
cargo run --release -- /path/to/data.csv

# Example
cargo run --release -- ../data/transactions_1m.csv
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

### Project Structure
```
src/
├── main.rs           # Entry point and CLI argument handling
├── process.rs        # Data processing and aggregation
└── benchmark.rs      # Benchmarking utilities
```

### Output
Results will be saved in the `../results/` directory with the following files:
- `rust_account_month.csv`
- `rust_department_month.csv`
- `rust_account_dept_month.csv`
- `rust_pivot_aggregation.csv` (if pivot option is enabled)

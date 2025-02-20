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
   cargo build
   ```

### Usage

#### Run with Default Settings
```bash
cargo run
```

#### Configuration
The behavior can be configured by modifying `main.rs`:
- File Paths: Modify the CSV file paths in the configuration section
- Aggregation Options: Enable/disable specific aggregation types

### Project Structure
```
src/
├── main.rs           # Entry point and configuration
├── process.rs        # Data processing and aggregation
└── benchmark.rs      # Benchmarking utilities
```

### Output
Results will be saved in the `../results/` directory with the following files:
- `rust_account_month.csv`
- `rust_department_month.csv`
- `rust_account_dept_month.csv`
- `rust_pivot_aggregation.csv` (if pivot option is enabled)

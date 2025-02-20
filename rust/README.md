# Rust Implementation for Accounting Data Aggregation Benchmark

## Overview

This Rust solution uses the [Polars](https://github.com/pola-rs/polars) library to:
- Generate a sample CSV dataset
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

- Data Generation: Uncomment the data generation function call to generate new sample data
- File Paths: Modify the CSV file paths in the configuration section
- Aggregation Options: Enable/disable specific aggregation types

### Project Structure
```
src/
├── main.rs           # Entry point and configuration
├── data_generator.rs # Sample data generation
├── process.rs        # Data processing and aggregation
└── benchmark.rs      # Benchmarking utilities
```

### Output
Results will be saved in the `../results/` directory with the following files:
- `rust_account_month.csv`
- `rust_department_month.csv`
- `rust_account_dept_month.csv`
- `rust_pivot_aggregation.csv`

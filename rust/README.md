# Rust Implementation for Accounting Data Aggregation Benchmark

## Overview
This Rust solution uses the Polars library to:
- Generate a sample CSV dataset.
- Perform three types of aggregations (Account × Month, Department × Month, Account × Department × Month).
- Benchmark the execution times.
- Optionally pivot the results for a pivot table format.

## Setup and Execution
1. Ensure you have Rust (LTS version) installed.
2. Navigate to the `rust` directory.
3. Build the project with:
cargo build

4. To run the benchmark (and optionally generate data), use:
cargo run

*Uncomment the data generation call in `main.rs` if you need to regenerate the CSV fil

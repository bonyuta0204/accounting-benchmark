# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a performance benchmarking project comparing Rust (Polars) and Go (Gota DataFrame) implementations for accounting data aggregation. The project processes CSV files with millions of rows and performs three types of aggregations:
1. Account × Monthly Aggregation
2. Department × Monthly Aggregation  
3. Account × Department × Monthly Aggregation

## Common Development Commands

### Data Generation
```bash
# Generate test data (1M, 10M, or 100M rows)
make generate-1m
make generate-10m
make generate-100m
make generate-all  # Generate all datasets
```

### Build Commands
```bash
# Build all components
make build-all

# Build specific components
make build-generator  # Build data generator (Rust)
make build-rust      # Build Rust benchmark implementation
make build-go        # Build Go benchmark implementation
```

### Running Benchmarks
```bash
# Run benchmarks for specific data sizes
make benchmark-1m    # Benchmark with 1M records
make benchmark-10m   # Benchmark with 10M records
make benchmark-100m  # Benchmark with 100M records
make benchmark-all   # Run all benchmarks

# Run individual implementations
cd rust && cargo run --release -- ../data/transactions_1m.csv
cd go && go run . ../data/transactions_1m.csv
```

### Cleanup
```bash
make clean-data  # Remove generated data files
make clean       # Clean all build artifacts and data
```

## Code Architecture

### Three Independent Components

1. **Generator** (`generator/`)
   - Rust CLI tool for generating test CSV data
   - Uses clap for command-line argument parsing
   - Configurable row count, date range, and output path

2. **Rust Implementation** (`rust/`)
   - `main.rs`: Entry point with CLI argument parsing (clap)
   - `benchmark.rs`: Orchestrates benchmarking and timing measurements
   - `process.rs`: Core data processing logic using Polars DataFrame
   - Always build with `--release` flag for accurate benchmarks

3. **Go Implementation** (`go/`)
   - `main.go`: Entry point and CLI argument handling
   - `benchmark.go`: Benchmark orchestration and timing
   - `process.go`: Core data processing using Gota DataFrame
   - No special build flags needed

### Data Flow Architecture

1. **Input**: CSV files with columns: Date, Account, Department, Amount
2. **Processing**: 
   - Parse date column to extract month
   - Group by required dimensions
   - Sum amounts
   - Optional pivot transformation
3. **Output**: Results saved to `results/` directory as CSV files

### Key Design Decisions

- Both implementations follow the same processing pattern for fair comparison
- Results are written to separate files with language-specific prefixes (rust_ or go_)
- Timing is measured for each aggregation type separately
- The Rust version shows ~9-14x performance improvement when built with --release
# Accounting Benchmark Makefile
# Benchmarks comparing Rust (Polars) and Go implementations for data aggregation

# Default shell
SHELL := /bin/bash

# Data sizes
SIZES := 1m 10m 100m

# Colors for help output
GREEN := \033[0;32m
YELLOW := \033[0;33m
CYAN := \033[0;36m
NC := \033[0m # No Color

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make ${CYAN}<target>${NC}\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  ${CYAN}%-20s${NC} %s\n", $$1, $$2 } /^##@/ { printf "\n${YELLOW}%s${NC}\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
all: build-all generate-all benchmark-all ## Build everything, generate all data, and run all benchmarks

##@ Build

.PHONY: build-generator
build-generator: ## Build the data generator (Rust)
	cd generator && cargo build --release

.PHONY: build-rust
build-rust: ## Build the Rust benchmark implementation
	cd rust && cargo build --release

.PHONY: build-go
build-go: ## Build the Go benchmark implementation
	cd go && go build

.PHONY: build-all
build-all: build-generator build-rust build-go ## Build all components

##@ Data Generation - CSV

.PHONY: generate-csv-1m
generate-csv-1m: build-generator ## Generate 1M rows CSV file
	@mkdir -p data
	cd generator && cargo run --release -- -o ../data/transactions_1m.csv -r 1000000

.PHONY: generate-csv-10m
generate-csv-10m: build-generator ## Generate 10M rows CSV file
	@mkdir -p data
	cd generator && cargo run --release -- -o ../data/transactions_10m.csv -r 10000000

.PHONY: generate-csv-100m
generate-csv-100m: build-generator ## Generate 100M rows CSV file
	@mkdir -p data
	cd generator && cargo run --release -- -o ../data/transactions_100m.csv -r 100000000

.PHONY: generate-csv-all
generate-csv-all: generate-csv-1m generate-csv-10m generate-csv-100m ## Generate all CSV files

##@ Data Generation - Parquet

.PHONY: generate-parquet-1m
generate-parquet-1m: build-generator ## Generate 1M rows Parquet file
	@mkdir -p data
	cd generator && cargo run --release -- -o ../data/transactions_1m.parquet -r 1000000 -f parquet

.PHONY: generate-parquet-10m
generate-parquet-10m: build-generator ## Generate 10M rows Parquet file
	@mkdir -p data
	cd generator && cargo run --release -- -o ../data/transactions_10m.parquet -r 10000000 -f parquet

.PHONY: generate-parquet-100m
generate-parquet-100m: build-generator ## Generate 100M rows Parquet file
	@mkdir -p data
	cd generator && cargo run --release -- -o ../data/transactions_100m.parquet -r 100000000 -f parquet

.PHONY: generate-parquet-all
generate-parquet-all: generate-parquet-1m generate-parquet-10m generate-parquet-100m ## Generate all Parquet files

.PHONY: generate-all
generate-all: generate-csv-all generate-parquet-all ## Generate all data files (CSV and Parquet)

# Legacy targets for backward compatibility
.PHONY: generate-1m generate-10m generate-100m
generate-1m: generate-csv-1m
generate-10m: generate-csv-10m
generate-100m: generate-csv-100m

##@ Benchmarks - Rust

.PHONY: rust-bench-csv-1m
rust-bench-csv-1m: build-rust data/transactions_1m.csv ## Run Rust benchmark with 1M CSV
	cd rust && cargo run --release -- ../data/transactions_1m.csv

.PHONY: rust-bench-csv-10m
rust-bench-csv-10m: build-rust data/transactions_10m.csv ## Run Rust benchmark with 10M CSV
	cd rust && cargo run --release -- ../data/transactions_10m.csv

.PHONY: rust-bench-csv-100m
rust-bench-csv-100m: build-rust data/transactions_100m.csv ## Run Rust benchmark with 100M CSV
	cd rust && cargo run --release -- ../data/transactions_100m.csv

.PHONY: rust-bench-parquet-1m
rust-bench-parquet-1m: build-rust data/transactions_1m.parquet ## Run Rust benchmark with 1M Parquet
	cd rust && cargo run --release -- ../data/transactions_1m.parquet

.PHONY: rust-bench-parquet-10m
rust-bench-parquet-10m: build-rust data/transactions_10m.parquet ## Run Rust benchmark with 10M Parquet
	cd rust && cargo run --release -- ../data/transactions_10m.parquet

.PHONY: rust-bench-parquet-100m
rust-bench-parquet-100m: build-rust data/transactions_100m.parquet ## Run Rust benchmark with 100M Parquet
	cd rust && cargo run --release -- ../data/transactions_100m.parquet

.PHONY: rust-bench-all
rust-bench-all: rust-bench-csv-1m rust-bench-csv-10m rust-bench-csv-100m rust-bench-parquet-1m rust-bench-parquet-10m rust-bench-parquet-100m ## Run all Rust benchmarks

##@ Benchmarks - Go

.PHONY: go-bench-csv-1m
go-bench-csv-1m: build-go data/transactions_1m.csv ## Run Go benchmark with 1M CSV
	cd go && ./accounting_benchmark_go ../data/transactions_1m.csv

.PHONY: go-bench-csv-10m
go-bench-csv-10m: build-go data/transactions_10m.csv ## Run Go benchmark with 10M CSV
	cd go && ./accounting_benchmark_go ../data/transactions_10m.csv

.PHONY: go-bench-csv-100m
go-bench-csv-100m: build-go data/transactions_100m.csv ## Run Go benchmark with 100M CSV
	cd go && ./accounting_benchmark_go ../data/transactions_100m.csv

.PHONY: go-bench-parquet-1m
go-bench-parquet-1m: build-go data/transactions_1m.parquet ## Run Go benchmark with 1M Parquet
	cd go && ./accounting_benchmark_go ../data/transactions_1m.parquet

.PHONY: go-bench-parquet-10m
go-bench-parquet-10m: build-go data/transactions_10m.parquet ## Run Go benchmark with 10M Parquet
	cd go && ./accounting_benchmark_go ../data/transactions_10m.parquet

.PHONY: go-bench-parquet-100m
go-bench-parquet-100m: build-go data/transactions_100m.parquet ## Run Go benchmark with 100M Parquet
	cd go && ./accounting_benchmark_go ../data/transactions_100m.parquet

.PHONY: go-bench-all
go-bench-all: go-bench-csv-1m go-bench-csv-10m go-bench-csv-100m go-bench-parquet-1m go-bench-parquet-10m go-bench-parquet-100m ## Run all Go benchmarks

##@ Combined Benchmarks

.PHONY: benchmark-csv-1m
benchmark-csv-1m: rust-bench-csv-1m go-bench-csv-1m ## Run both Rust and Go benchmarks with 1M CSV

.PHONY: benchmark-csv-10m
benchmark-csv-10m: rust-bench-csv-10m go-bench-csv-10m ## Run both Rust and Go benchmarks with 10M CSV

.PHONY: benchmark-csv-100m
benchmark-csv-100m: rust-bench-csv-100m go-bench-csv-100m ## Run both Rust and Go benchmarks with 100M CSV

.PHONY: benchmark-parquet-1m
benchmark-parquet-1m: rust-bench-parquet-1m go-bench-parquet-1m ## Run both Rust and Go benchmarks with 1M Parquet

.PHONY: benchmark-parquet-10m
benchmark-parquet-10m: rust-bench-parquet-10m go-bench-parquet-10m ## Run both Rust and Go benchmarks with 10M Parquet

.PHONY: benchmark-parquet-100m
benchmark-parquet-100m: rust-bench-parquet-100m go-bench-parquet-100m ## Run both Rust and Go benchmarks with 100M Parquet

.PHONY: benchmark-all
benchmark-all: rust-bench-all go-bench-all ## Run all benchmarks (Rust and Go, CSV and Parquet)

# Legacy targets for backward compatibility
.PHONY: benchmark-1m benchmark-10m benchmark-100m rust-bench-1m rust-bench-10m rust-bench-100m go-bench-1m go-bench-10m go-bench-100m
benchmark-1m: benchmark-csv-1m
benchmark-10m: benchmark-csv-10m
benchmark-100m: benchmark-csv-100m
rust-bench-1m: rust-bench-csv-1m
rust-bench-10m: rust-bench-csv-10m
rust-bench-100m: rust-bench-csv-100m
go-bench-1m: go-bench-csv-1m
go-bench-10m: go-bench-csv-10m
go-bench-100m: go-bench-csv-100m

##@ Comparison Reports

.PHONY: compare-formats
compare-formats: ## Compare CSV vs Parquet performance (requires 1M data files)
	@echo "===== Format Comparison Report ====="
	@echo ""
	@echo "CSV Performance (1M rows):"
	@$(MAKE) -s go-bench-csv-1m | grep -E "(Summary for|Load time:|Total time:)" || true
	@echo ""
	@echo "Parquet Performance (1M rows):"
	@$(MAKE) -s go-bench-parquet-1m | grep -E "(Summary for|Load time:|Total time:)" || true

.PHONY: compare-languages
compare-languages: ## Compare Rust vs Go performance (requires 1M CSV)
	@echo "===== Language Comparison Report ====="
	@echo ""
	@echo "Rust Performance (1M CSV):"
	@$(MAKE) -s rust-bench-csv-1m || true
	@echo ""
	@echo "Go Performance (1M CSV):"
	@$(MAKE) -s go-bench-csv-1m | grep -E "(Summary for|Total time:)" || true

##@ Data Dependencies

# Auto-generate data files if they don't exist
data/transactions_1m.csv:
	$(MAKE) generate-csv-1m

data/transactions_10m.csv:
	$(MAKE) generate-csv-10m

data/transactions_100m.csv:
	$(MAKE) generate-csv-100m

data/transactions_1m.parquet:
	$(MAKE) generate-parquet-1m

data/transactions_10m.parquet:
	$(MAKE) generate-parquet-10m

data/transactions_100m.parquet:
	$(MAKE) generate-parquet-100m

##@ Cleanup

.PHONY: clean
clean: ## Clean all build artifacts
	cd generator && cargo clean
	cd rust && cargo clean
	cd go && go clean
	rm -f go/accounting_benchmark_go

.PHONY: clean-data
clean-data: ## Remove all generated data files
	rm -f data/transactions_*.csv data/transactions_*.parquet

.PHONY: clean-results
clean-results: ## Remove all benchmark result files
	rm -f results/*.csv
	rm -f data/*_*.csv  # Remove result files in data directory

.PHONY: clean-all
clean-all: clean clean-data clean-results ## Clean everything (build artifacts, data, and results)

##@ Quick Commands

.PHONY: quick-test
quick-test: generate-csv-1m generate-parquet-1m ## Quick test with 1M rows (CSV and Parquet)
	@echo "Running quick benchmark test..."
	@$(MAKE) compare-formats

.PHONY: setup
setup: build-all ## Initial setup - build all components
	@echo "Setup complete! You can now run 'make generate-all' to create test data."
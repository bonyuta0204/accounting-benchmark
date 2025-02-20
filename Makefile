.PHONY: all clean clean-data benchmark-all generate-all build-all

# Build targets
build-generator:
	cd generator && cargo build --release

build-rust:
	cd rust && cargo build --release

build-go:
	cd go && go build

build-all: build-generator build-rust build-go

# Data generation targets
generate-1m: build-generator
	cd generator && cargo run --release -- -o ../data/transactions_1m.csv -r 1000000

generate-10m: build-generator
	cd generator && cargo run --release -- -o ../data/transactions_10m.csv -r 10000000

generate-100m: build-generator
	cd generator && cargo run --release -- -o ../data/transactions_100m.csv -r 100000000

generate-all: generate-1m generate-10m generate-100m

# Rust benchmark targets
rust-bench-1m: build-rust data/transactions_1m.csv
	cd rust && cargo run --release -- ../data/transactions_1m.csv

rust-bench-10m: build-rust data/transactions_10m.csv
	cd rust && cargo run --release -- ../data/transactions_10m.csv

rust-bench-100m: build-rust data/transactions_100m.csv
	cd rust && cargo run --release -- ../data/transactions_100m.csv

# Go benchmark targets
go-bench-1m: build-go data/transactions_1m.csv
	cd go && go run . ../data/transactions_1m.csv

go-bench-10m: build-go data/transactions_10m.csv
	cd go && go run . ../data/transactions_10m.csv

go-bench-100m: build-go data/transactions_100m.csv
	cd go && go run . ../data/transactions_100m.csv

# Data file dependencies
data/transactions_1m.csv:
	$(MAKE) generate-1m

data/transactions_10m.csv:
	$(MAKE) generate-10m

data/transactions_100m.csv:
	$(MAKE) generate-100m

# Clean targets
clean:
	cd generator && cargo clean
	cd rust && cargo clean
	cd go && go clean
	
clean-data:
	rm -f data/transactions_*.csv

# Combined benchmark targets
benchmark-1m: rust-bench-1m go-bench-1m

benchmark-10m: rust-bench-10m go-bench-10m

benchmark-100m: rust-bench-100m go-bench-100m

benchmark-all: rust-bench-1m rust-bench-10m rust-bench-100m go-bench-1m go-bench-10m go-bench-100m

# Default target
all: build-all generate-all benchmark-all

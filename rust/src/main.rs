mod benchmark;
mod process;

fn main() {
    // Path to the CSV data (adjust if necessary)
    let csv_path = "../data/sample_1m.csv";

    // Run the aggregation benchmarks
    benchmark::run_benchmarks(csv_path);
}

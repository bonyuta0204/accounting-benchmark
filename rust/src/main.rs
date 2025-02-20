mod benchmark;
mod data_generator;
mod process;

fn main() {
    // Path to the CSV data (adjust if necessary)
    let csv_path = "../data/sample_1m.csv";

    // Uncomment the next line to generate new data:
    data_generator::generate_csv(csv_path, 1_000_000).expect("Failed to generate CSV");

    // Run the aggregation benchmarks
    benchmark::run_benchmarks(csv_path);
}

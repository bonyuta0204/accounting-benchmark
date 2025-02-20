mod benchmark;
mod process;
use clap::Parser;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Path to the CSV file
    csv_path: String,
}

fn main() {
    let args = Args::parse();
    benchmark::run_benchmarks(&args.csv_path);
}

mod aggregator;
mod benchmark;
mod benchmark_framework;
mod process;

use aggregator::get_all_aggregators;
use benchmark_framework::run_comparison;
use clap::Parser;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Path to the CSV file
    csv_path: String,
    
    /// Use legacy benchmark mode
    #[arg(long, default_value_t = false)]
    legacy: bool,
}

fn main() {
    let args = Args::parse();
    
    if args.legacy {
        // Use old benchmark for compatibility
        benchmark::run_benchmarks(&args.csv_path);
    } else {
        // Use new framework
        let aggregators = get_all_aggregators();
        if let Err(e) = run_comparison(&args.csv_path, aggregators) {
            eprintln!("Benchmark failed: {}", e);
            std::process::exit(1);
        }
    }
}

use crate::aggregator::Aggregator;
use crate::process::{add_month_column, load_csv};
use polars::prelude::*;
use std::error::Error;
use std::fs::File;
use std::time::Instant;

pub struct BenchmarkResult {
    pub name: String,
    pub total_time_ms: f64,
    pub processing_time_ms: f64,
    pub io_time_ms: f64,
}

impl BenchmarkResult {
    pub fn print(&self) {
        println!("\n=== {} ===", self.name);
        println!("Total time (with I/O): {:.3}ms", self.total_time_ms);
        println!("Processing time (without I/O): {:.3}ms", self.processing_time_ms);
        println!("I/O overhead: {:.3}ms", self.io_time_ms);
    }
}

pub struct BenchmarkFramework {
    csv_path: String,
    aggregators: Vec<Box<dyn Aggregator>>,
}

impl BenchmarkFramework {
    pub fn new(csv_path: String) -> Self {
        Self {
            csv_path,
            aggregators: Vec::new(),
        }
    }

    pub fn add_aggregator(mut self, aggregator: Box<dyn Aggregator>) -> Self {
        self.aggregators.push(aggregator);
        self
    }

    pub fn add_aggregators(mut self, aggregators: Vec<Box<dyn Aggregator>>) -> Self {
        self.aggregators.extend(aggregators);
        self
    }

    pub fn run(&self) -> Result<Vec<BenchmarkResult>, Box<dyn Error>> {
        // First load the data once
        println!("Loading CSV data...");
        let load_start = Instant::now();
        let df = load_csv(&self.csv_path)?;
        let load_duration = load_start.elapsed();
        println!("CSV loading took: {:?}", load_duration);
        println!("Loaded {} rows", df.height());
        
        // Add month column once
        let df_with_month = add_month_column(df)?;
        
        let mut results = Vec::new();
        
        for aggregator in &self.aggregators {
            // Measure total time (including I/O)
            let total_start = Instant::now();
            let result_with_io = self.run_aggregation_with_io(aggregator.as_ref())?;
            let total_duration = total_start.elapsed();
            
            // Measure processing-only time
            let proc_start = Instant::now();
            let result_without_io = aggregator.aggregate(&df_with_month)?;
            let proc_duration = proc_start.elapsed();
            
            // Write results
            self.write_results(&result_with_io, aggregator.output_file())?;
            
            let benchmark_result = BenchmarkResult {
                name: aggregator.name().to_string(),
                total_time_ms: total_duration.as_secs_f64() * 1000.0,
                processing_time_ms: proc_duration.as_secs_f64() * 1000.0,
                io_time_ms: (total_duration - proc_duration).as_secs_f64() * 1000.0,
            };
            
            benchmark_result.print();
            results.push(benchmark_result);
        }
        
        Ok(results)
    }

    fn run_aggregation_with_io(&self, aggregator: &dyn Aggregator) -> Result<DataFrame, Box<dyn Error>> {
        // This includes file I/O
        let df = load_csv(&self.csv_path)?;
        let df_with_month = add_month_column(df)?;
        aggregator.aggregate(&df_with_month)
    }

    fn write_results(&self, df: &DataFrame, path: &str) -> Result<(), Box<dyn Error>> {
        let mut file = File::create(path)?;
        let mut df = df.clone();
        CsvWriter::new(&mut file).finish(&mut df)?;
        Ok(())
    }
}

/// Run comparison between multiple aggregators
pub fn run_comparison(csv_path: &str, aggregators: Vec<Box<dyn Aggregator>>) -> Result<(), Box<dyn Error>> {
    let framework = BenchmarkFramework::new(csv_path.to_string())
        .add_aggregators(aggregators);
    
    let results = framework.run()?;
    
    // Print summary
    println!("\n=== BENCHMARK SUMMARY ===");
    println!("{:<50} {:>15} {:>15} {:>15}", "Aggregation", "Total (ms)", "Processing (ms)", "I/O (ms)");
    println!("{:-<50} {:-^15} {:-^15} {:-^15}", "", "", "", "");
    
    for result in &results {
        println!("{:<50} {:>15.3} {:>15.3} {:>15.3}",
            result.name,
            result.total_time_ms,
            result.processing_time_ms,
            result.io_time_ms
        );
    }
    
    // Find fastest/slowest
    if let (Some(fastest), Some(slowest)) = (
        results.iter().min_by(|a, b| a.processing_time_ms.partial_cmp(&b.processing_time_ms).unwrap()),
        results.iter().max_by(|a, b| a.processing_time_ms.partial_cmp(&b.processing_time_ms).unwrap())
    ) {
        println!("\nFastest: {} ({:.3}ms)", fastest.name, fastest.processing_time_ms);
        println!("Slowest: {} ({:.3}ms)", slowest.name, slowest.processing_time_ms);
        println!("Speed difference: {:.2}x", slowest.processing_time_ms / fastest.processing_time_ms);
    }
    
    Ok(())
}
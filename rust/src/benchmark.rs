use crate::process::{
    aggregate_by_account_department_month, aggregate_by_account_month,
    aggregate_by_department_month, pivot_aggregation,
    load_csv, aggregate_by_account_month_df, aggregate_by_department_month_df,
    aggregate_by_account_department_month_df,
};
use polars::prelude::*;
use std::error::Error;
use std::fs::File;
use std::time::Instant;

pub fn run_benchmarks(csv_path: &str) {
    // First load the data once
    println!("Loading CSV data...");
    let load_start = Instant::now();
    let df = load_csv(csv_path).expect("Failed to load CSV");
    let load_duration = load_start.elapsed();
    println!("CSV loading took: {:?}", load_duration);
    println!("Loaded {} rows", df.height());
    println!();

    // Benchmark Account × Monthly Aggregation
    println!("=== Account × Monthly Aggregation ===");
    // Total time (including I/O)
    let start = Instant::now();
    let account_month_df = aggregate_by_account_month(csv_path).expect("Aggregation failed");
    let total_duration = start.elapsed();
    println!("Total time (with I/O): {:?}", total_duration);
    
    // Processing-only time
    let start = Instant::now();
    let account_month_df_proc = aggregate_by_account_month_df(&df).expect("Aggregation failed");
    let proc_duration = start.elapsed();
    println!("Processing time (without I/O): {:?}", proc_duration);
    println!("I/O overhead: {:?}", total_duration - proc_duration);
    
    write_df_to_csv(&account_month_df, "../results/rust_account_month.csv").unwrap();

    // Benchmark Department × Monthly Aggregation
    println!("\n=== Department × Monthly Aggregation ===");
    // Total time (including I/O)
    let start = Instant::now();
    let department_month_df = aggregate_by_department_month(csv_path).expect("Aggregation failed");
    let total_duration = start.elapsed();
    println!("Total time (with I/O): {:?}", total_duration);
    
    // Processing-only time
    let start = Instant::now();
    let department_month_df_proc = aggregate_by_department_month_df(&df).expect("Aggregation failed");
    let proc_duration = start.elapsed();
    println!("Processing time (without I/O): {:?}", proc_duration);
    println!("I/O overhead: {:?}", total_duration - proc_duration);
    
    write_df_to_csv(&department_month_df, "../results/rust_department_month.csv").unwrap();

    // Benchmark Account × Department × Monthly Aggregation
    println!("\n=== Account × Department × Monthly Aggregation ===");
    // Total time (including I/O)
    let start = Instant::now();
    let account_dept_month_df =
        aggregate_by_account_department_month(csv_path).expect("Aggregation failed");
    let total_duration = start.elapsed();
    println!("Total time (with I/O): {:?}", total_duration);
    
    // Processing-only time
    let start = Instant::now();
    let account_dept_month_df_proc = aggregate_by_account_department_month_df(&df).expect("Aggregation failed");
    let proc_duration = start.elapsed();
    println!("Processing time (without I/O): {:?}", proc_duration);
    println!("I/O overhead: {:?}", total_duration - proc_duration);
    
    write_df_to_csv(
        &account_dept_month_df,
        "../results/rust_account_dept_month.csv",
    )
    .unwrap();

    // Optional: Pivot aggregation result (using account-month aggregation as example)
    let pivot_df = pivot_aggregation(&account_month_df, &["Account"]).expect("Pivot failed");
    write_df_to_csv(&pivot_df, "../results/rust_pivot_aggregation.csv").unwrap();
}

fn write_df_to_csv(df: &DataFrame, path: &str) -> Result<(), Box<dyn Error>> {
    let mut file = File::create(path)?;
    let mut df = df.clone();
    CsvWriter::new(&mut file).finish(&mut df)?;
    Ok(())
}

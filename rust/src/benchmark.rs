use std::time::Instant;
use crate::process::{aggregate_by_account_month, aggregate_by_department_month, aggregate_by_account_department_month, pivot_aggregation};
use polars::prelude::*;
use std::error::Error;
use std::fs::File;

pub fn run_benchmarks(csv_path: &str) {
    // Benchmark Account × Monthly Aggregation
    let start = Instant::now();
    let account_month_df = aggregate_by_account_month(csv_path).expect("Aggregation failed");
    let duration = start.elapsed();
    println!("Account × Monthly Aggregation took: {:?}", duration);
    write_df_to_csv(&account_month_df, "../results/rust_account_month.csv").unwrap();

    // Benchmark Department × Monthly Aggregation
    let start = Instant::now();
    let department_month_df = aggregate_by_department_month(csv_path).expect("Aggregation failed");
    let duration = start.elapsed();
    println!("Department × Monthly Aggregation took: {:?}", duration);
    write_df_to_csv(&department_month_df, "../results/rust_department_month.csv").unwrap();

    // Benchmark Account × Department × Monthly Aggregation
    let start = Instant::now();
    let account_dept_month_df = aggregate_by_account_department_month(csv_path).expect("Aggregation failed");
    let duration = start.elapsed();
    println!("Account × Department × Monthly Aggregation took: {:?}", duration);
    write_df_to_csv(&account_dept_month_df, "../results/rust_account_dept_month.csv").unwrap();

    // Optional: Pivot aggregation result (using account-month aggregation as example)
    let pivot_df = pivot_aggregation(&account_month_df, &["Account"]).expect("Pivot failed");
    write_df_to_csv(&pivot_df, "../results/rust_pivot_aggregation.csv").unwrap();
}

fn write_df_to_csv(df: &DataFrame, path: &str) -> Result<(), Box<dyn Error>> {
    let mut file = File::create(path)?;
    let mut df = df.clone();
    CsvWriter::new(&mut file)
        .finish(&mut df)?;
    Ok(())
}

use chrono::{Duration, NaiveDate};
use clap::{Parser, ValueEnum};
use polars::prelude::*;
use rand::Rng;
use std::error::Error;
use std::fs::File;
use std::io::Write;
use std::path::PathBuf;

#[derive(ValueEnum, Clone, Debug)]
enum Format {
    Csv,
    Parquet,
}

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Output file path
    #[arg(short, long)]
    output: PathBuf,

    /// Number of rows to generate
    #[arg(short, long, default_value_t = 1000)]
    rows: usize,

    /// Start date (YYYY-MM-DD)
    #[arg(short, long, default_value = "2020-01-01")]
    start_date: String,

    /// End date (YYYY-MM-DD)
    #[arg(short, long, default_value = "2020-12-31")]
    end_date: String,
    
    /// Output format
    #[arg(short, long, value_enum, default_value_t = Format::Csv)]
    format: Format,
}

fn generate_csv(
    path: PathBuf,
    rows: usize,
    start_date: NaiveDate,
    end_date: NaiveDate,
) -> Result<(), Box<dyn Error>> {
    let mut file = File::create(path)?;
    writeln!(file, "Date,Amount,Account,Department")?;

    let date_range = end_date.signed_duration_since(start_date).num_days() as i64;
    let accounts = vec!["Sales", "Expenses", "Assets", "Liabilities"];
    let departments = vec!["Sales", "Development", "HR", "Finance"];

    let mut rng = rand::thread_rng();

    for _ in 0..rows {
        let offset = rng.gen_range(0..=date_range);
        let date = start_date + Duration::days(offset);
        let amount: f64 = rng.gen_range(10.0..1000.0);
        let account = accounts[rng.gen_range(0..accounts.len())];
        let department = departments[rng.gen_range(0..departments.len())];

        writeln!(
            file,
            "{},{:.2},{},{}",
            date.format("%Y-%m-%d"),
            amount,
            account,
            department
        )?;
    }
    Ok(())
}

fn generate_parquet(
    path: PathBuf,
    rows: usize,
    start_date: NaiveDate,
    end_date: NaiveDate,
) -> Result<(), Box<dyn Error>> {
    let date_range = end_date.signed_duration_since(start_date).num_days() as i64;
    let accounts = vec!["Sales", "Expenses", "Assets", "Liabilities"];
    let departments = vec!["Sales", "Development", "HR", "Finance"];

    let mut rng = rand::thread_rng();
    
    let mut dates = Vec::with_capacity(rows);
    let mut amounts = Vec::with_capacity(rows);
    let mut account_vec = Vec::with_capacity(rows);
    let mut department_vec = Vec::with_capacity(rows);

    for _ in 0..rows {
        let offset = rng.gen_range(0..=date_range);
        let date = start_date + Duration::days(offset);
        dates.push(date.format("%Y-%m-%d").to_string());
        
        let amount: f64 = rng.gen_range(10.0..1000.0);
        amounts.push(amount);
        
        let account = accounts[rng.gen_range(0..accounts.len())];
        account_vec.push(account.to_string());
        
        let department = departments[rng.gen_range(0..departments.len())];
        department_vec.push(department.to_string());
    }

    let df = DataFrame::new(vec![
        Series::new("Date", dates),
        Series::new("Amount", amounts),
        Series::new("Account", account_vec),
        Series::new("Department", department_vec),
    ])?;

    let mut file = File::create(&path)?;
    ParquetWriter::new(&mut file).finish(&mut df.clone())?;
    
    Ok(())
}

fn main() -> Result<(), Box<dyn Error>> {
    let args = Args::parse();

    let start_date = NaiveDate::parse_from_str(&args.start_date, "%Y-%m-%d")?;
    let end_date = NaiveDate::parse_from_str(&args.end_date, "%Y-%m-%d")?;

    if end_date <= start_date {
        return Err("End date must be after start date".into());
    }

    match args.format {
        Format::Csv => generate_csv(args.output, args.rows, start_date, end_date)?,
        Format::Parquet => generate_parquet(args.output, args.rows, start_date, end_date)?,
    }
    
    println!("Successfully generated {} rows of data in {:?} format", args.rows, args.format);
    Ok(())
}

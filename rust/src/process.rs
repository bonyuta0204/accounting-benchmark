use polars::prelude::*;
use std::error::Error;
use std::fs::File;

pub fn aggregate_by_account_month(csv_path: &str) -> Result<DataFrame, Box<dyn Error>> {
    // Read CSV file
    let file = File::open(csv_path)?;
    let mut df = CsvReader::new(file).finish()?;

    // Create a new column "Month" extracted from the "Date" column
    df = df
        .lazy()
        .with_column(col("Date").dt().month().alias("Month"))
        .collect()?;

    // Group by Account and Month; sum the Amount
    let agg_df = df
        .lazy()
        .group_by(&[col("Account"), col("Month")])
        .agg([col("Amount").sum().alias("Total")])
        .collect()?;

    Ok(agg_df)
}

pub fn aggregate_by_department_month(csv_path: &str) -> Result<DataFrame, Box<dyn Error>> {
    let file = File::open(csv_path)?;
    let mut df = CsvReader::new(file).finish()?;

    df = df
        .lazy()
        .with_column(col("Date").dt().month().alias("Month"))
        .collect()?;

    let agg_df = df
        .lazy()
        .group_by(&[col("Department"), col("Month")])
        .agg([col("Amount").sum().alias("Total")])
        .collect()?;

    Ok(agg_df)
}

pub fn aggregate_by_account_department_month(csv_path: &str) -> Result<DataFrame, Box<dyn Error>> {
    let file = File::open(csv_path)?;
    let mut df = CsvReader::new(file).finish()?;

    df = df
        .lazy()
        .with_column(col("Date").dt().month().alias("Month"))
        .collect()?;

    let agg_df = df
        .lazy()
        .group_by(&[col("Account"), col("Department"), col("Month")])
        .agg([col("Amount").sum().alias("Total")])
        .collect()?;

    Ok(agg_df)
}

// Optional: Pivot the aggregation result (e.g. Account × Month where each unique month becomes a column)
pub fn pivot_aggregation(df: &DataFrame, group_cols: &[&str]) -> Result<DataFrame, Box<dyn Error>> {
    // For pivot, we need to use a different approach with dynamic groupby
    let months = df.column("Month")?.unique()?;
    let mut agg_exprs = Vec::new();

    // Convert months to i32 for filtering
    for month in months.i32()? {
        if let Some(m) = month {
            let month_expr = col("Total")
                .filter(col("Month").eq(lit(m)))
                .sum()
                .alias(&format!("Month_{}", m));
            agg_exprs.push(month_expr);
        }
    }

    let pivot_df = df
        .clone()
        .lazy()
        .group_by(group_cols)
        .agg(&agg_exprs)
        .collect()?;

    Ok(pivot_df)
}

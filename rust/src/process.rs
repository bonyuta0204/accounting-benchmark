use polars::prelude::*;
use std::error::Error;

pub fn aggregate_by_account_month(csv_path: &str) -> Result<DataFrame, Box<dyn Error>> {
    // Read CSV file
    let mut df = CsvReader::from_path(csv_path)?
        .infer_schema(Some(100))
        .has_header(true)
        .finish()?;
    
    // Create a new column "Month" extracted from the "Date" column
    df = df.lazy()
        .with_column(
            col("Date")
            .str()
            .strptime(DataType::Date, StrptimeOptions {
                format: Some("%Y-%m-%d".into()),
                strict: false,
                exact: true,
                cache: true,
            })
            .dt()
            .month()
            .alias("Month")
        )
        .collect()?;

    // Group by Account and Month; sum the Amount
    let agg_df = df.lazy()
        .groupby(&[col("Account"), col("Month")])
        .agg([col("Amount").sum().alias("Total")])
        .sort("Account", Default::default())
        .collect()?;

    Ok(agg_df)
}

pub fn aggregate_by_department_month(csv_path: &str) -> Result<DataFrame, Box<dyn Error>> {
    let mut df = CsvReader::from_path(csv_path)?
        .infer_schema(Some(100))
        .has_header(true)
        .finish()?;

    df = df.lazy()
        .with_column(
            col("Date")
            .str()
            .strptime(DataType::Date, StrptimeOptions {
                format: Some("%Y-%m-%d".into()),
                strict: false,
                exact: true,
                cache: true,
            })
            .dt()
            .month()
            .alias("Month")
        )
        .collect()?;

    let agg_df = df.lazy()
        .groupby(&[col("Department"), col("Month")])
        .agg([col("Amount").sum().alias("Total")])
        .sort("Department", Default::default())
        .collect()?;

    Ok(agg_df)
}

pub fn aggregate_by_account_department_month(csv_path: &str) -> Result<DataFrame, Box<dyn Error>> {
    let mut df = CsvReader::from_path(csv_path)?
        .infer_schema(Some(100))
        .has_header(true)
        .finish()?;

    df = df.lazy()
        .with_column(
            col("Date")
            .str()
            .strptime(DataType::Date, StrptimeOptions {
                format: Some("%Y-%m-%d".into()),
                strict: false,
                exact: true,
                cache: true,
            })
            .dt()
            .month()
            .alias("Month")
        )
        .collect()?;

    let agg_df = df.lazy()
        .groupby(&[col("Account"), col("Department"), col("Month")])
        .agg([col("Amount").sum().alias("Total")])
        .sort("Account", Default::default())
        .collect()?;

    Ok(agg_df)
}

// Optional: Pivot the aggregation result (e.g. Account × Month where each unique month becomes a column)
pub fn pivot_aggregation(df: &DataFrame, group_cols: &[&str]) -> Result<DataFrame, Box<dyn Error>> {
    let pivot_df = df.pivot(
        group_cols,
        "Month",
        "Total",
        PivotAgg::Sum,
    )?;
    Ok(pivot_df)
}


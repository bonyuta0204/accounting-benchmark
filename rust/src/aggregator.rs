use polars::prelude::*;
use std::error::Error;

/// Trait for defining aggregation operations
pub trait Aggregator: Send + Sync {
    /// Name of the aggregation for display
    fn name(&self) -> &str;
    
    /// Output file name for the results
    fn output_file(&self) -> &str;
    
    /// Perform the aggregation on the given DataFrame
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>>;
}

/// Account × Month aggregation
pub struct AccountMonthAggregator;

impl Aggregator for AccountMonthAggregator {
    fn name(&self) -> &str {
        "Account × Monthly Aggregation"
    }
    
    fn output_file(&self) -> &str {
        "../results/rust_account_month.csv"
    }
    
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>> {
        df.clone()
            .lazy()
            .group_by(&[col("Account"), col("Month")])
            .agg([col("Amount").sum().alias("Total")])
            .collect()
            .map_err(|e| e.into())
    }
}

/// Department × Month aggregation
pub struct DepartmentMonthAggregator;

impl Aggregator for DepartmentMonthAggregator {
    fn name(&self) -> &str {
        "Department × Monthly Aggregation"
    }
    
    fn output_file(&self) -> &str {
        "../results/rust_department_month.csv"
    }
    
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>> {
        df.clone()
            .lazy()
            .group_by(&[col("Department"), col("Month")])
            .agg([col("Amount").sum().alias("Total")])
            .collect()
            .map_err(|e| e.into())
    }
}

/// Account × Department × Month aggregation
pub struct AccountDepartmentMonthAggregator;

impl Aggregator for AccountDepartmentMonthAggregator {
    fn name(&self) -> &str {
        "Account × Department × Monthly Aggregation"
    }
    
    fn output_file(&self) -> &str {
        "../results/rust_account_dept_month.csv"
    }
    
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>> {
        df.clone()
            .lazy()
            .group_by(&[col("Account"), col("Department"), col("Month")])
            .agg([col("Amount").sum().alias("Total")])
            .collect()
            .map_err(|e| e.into())
    }
}

/// Monthly total aggregation
pub struct MonthlyTotalAggregator;

impl Aggregator for MonthlyTotalAggregator {
    fn name(&self) -> &str {
        "Monthly Total Aggregation"
    }
    
    fn output_file(&self) -> &str {
        "../results/rust_monthly_total.csv"
    }
    
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>> {
        df.clone()
            .lazy()
            .group_by(&[col("Month")])
            .agg([col("Amount").sum().alias("Total")])
            .collect()
            .map_err(|e| e.into())
    }
}

/// Account total aggregation
pub struct AccountTotalAggregator;

impl Aggregator for AccountTotalAggregator {
    fn name(&self) -> &str {
        "Account Total Aggregation"
    }
    
    fn output_file(&self) -> &str {
        "../results/rust_account_total.csv"
    }
    
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>> {
        df.clone()
            .lazy()
            .group_by(&[col("Account")])
            .agg([col("Amount").sum().alias("Total")])
            .collect()
            .map_err(|e| e.into())
    }
}

/// Department total aggregation
pub struct DepartmentTotalAggregator;

impl Aggregator for DepartmentTotalAggregator {
    fn name(&self) -> &str {
        "Department Total Aggregation"
    }
    
    fn output_file(&self) -> &str {
        "../results/rust_department_total.csv"
    }
    
    fn aggregate(&self, df: &DataFrame) -> Result<DataFrame, Box<dyn Error>> {
        df.clone()
            .lazy()
            .group_by(&[col("Department")])
            .agg([col("Amount").sum().alias("Total")])
            .collect()
            .map_err(|e| e.into())
    }
}

/// Get all available aggregators
pub fn get_all_aggregators() -> Vec<Box<dyn Aggregator>> {
    vec![
        Box::new(AccountMonthAggregator),
        Box::new(DepartmentMonthAggregator),
        Box::new(AccountDepartmentMonthAggregator),
        Box::new(MonthlyTotalAggregator),
        Box::new(AccountTotalAggregator),
        Box::new(DepartmentTotalAggregator),
    ]
}
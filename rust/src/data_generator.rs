use chrono::{Duration, NaiveDate};
use rand::Rng;
use std::error::Error;
use std::fs::File;
use std::io::Write;

pub fn generate_csv(path: &str, rows: usize) -> Result<(), Box<dyn Error>> {
    let mut file = File::create(path)?;
    writeln!(file, "Date,Amount,Account,Department")?;

    let start_date = NaiveDate::from_ymd(2020, 1, 1);
    let accounts = vec!["Sales", "Expenses", "Assets", "Liabilities"];
    let departments = vec!["Sales", "Development", "HR", "Finance"];

    let mut rng = rand::thread_rng();

    for _ in 0..rows {
        let offset = rng.gen_range(0..365);
        let date = start_date + Duration::days(offset);
        let amount: f64 = rng.gen_range(10.0..1000.0);
        let account = accounts[rng.gen_range(0..accounts.len())];
        let department = departments[rng.gen_range(0..departments.len())];

        writeln!(
            file,
            "{},{},{},{}",
            date.format("%Y-%m-%d"),
            amount,
            account,
            department
        )?;
    }
    Ok(())
}

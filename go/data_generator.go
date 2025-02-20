package main

import (
    "encoding/csv"
    "fmt"
    "math/rand"
    "os"
    "time"
)

func GenerateCSV(path string, rows int) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write header
    header := []string{"Date", "Amount", "Account", "Department"}
    writer.Write(header)

    startDate, _ := time.Parse("2006-01-02", "2020-01-01")
    accounts := []string{"Sales", "Expenses", "Assets", "Liabilities"}
    departments := []string{"Sales", "Development", "HR", "Finance"}

    rand.Seed(time.Now().UnixNano())
    for i := 0; i < rows; i++ {
        offset := rand.Intn(365)
        date := startDate.AddDate(0, 0, offset).Format("2006-01-02")
        amount := rand.Float64()*990 + 10 // random float between 10 and 1000
        account := accounts[rand.Intn(len(accounts))]
        department := departments[rand.Intn(len(departments))]

        record := []string{
            date,
            fmt.Sprintf("%.2f", amount),
            account,
            department,
        }
        writer.Write(record)
    }
    return nil
}

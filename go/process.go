package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/go-gota/gota/dataframe"
    "github.com/go-gota/gota/series"
)

// ProcessAggregations performs the aggregations and writes output CSV files.
func ProcessAggregations(csvPath string) {
    // Open CSV file
    file, err := os.Open(csvPath)
    if err != nil {
        fmt.Println("Error opening CSV:", err)
        return
    }
    defer file.Close()

    df := dataframe.ReadCSV(file)

    // Add a Month column extracted from the Date column
    dates := df.Col("Date").Records()[1:] // skip header row if present
    months := make([]int, len(dates))
    for i, d := range dates {
        t, err := time.Parse("2006-01-02", d)
        if err != nil {
            months[i] = 0
        } else {
            months[i] = int(t.Month())
        }
    }
    monthSeries := series.Ints(months).Name("Month")
    df = df.Mutate(monthSeries)

    // Aggregation 1: Account × Monthly Aggregation
    accountAgg := aggregate(df, "Account", "Month")
    writeDataFrameToCSV(accountAgg, "../results/go_account_month.csv")

    // Aggregation 2: Department × Monthly Aggregation
    departmentAgg := aggregate(df, "Department", "Month")
    writeDataFrameToCSV(departmentAgg, "../results/go_department_month.csv")

    // Aggregation 3: Account × Department × Monthly Aggregation
    accountDeptAgg := aggregateTwo(df, "Account", "Department", "Month")
    writeDataFrameToCSV(accountDeptAgg, "../results/go_account_dept_month.csv")

    // Optionally: Create a pivot table for account aggregation (rows: Account, columns: Month)
    pivot := pivotAggregation(accountAgg, "Account", "Month", "Total")
    writeDataFrameToCSV(pivot, "../results/go_pivot_aggregation.csv")
}

func aggregate(df dataframe.DataFrame, groupCol string, monthCol string) dataframe.DataFrame {
    // Group by groupCol and monthCol, summing the "Amount" column.
    agg := df.GroupBy(groupCol, monthCol).Aggregation([]dataframe.AggregationOptions{
        {Column: "Amount", Func: dataframe.Sum},
    })
    // Rename the aggregated column to "Total"
    agg = agg.Rename("Amount_sum", "Total")
    return agg
}

func aggregateTwo(df dataframe.DataFrame, groupCol1, groupCol2, monthCol string) dataframe.DataFrame {
    agg := df.GroupBy(groupCol1, groupCol2, monthCol).Aggregation([]dataframe.AggregationOptions{
        {Column: "Amount", Func: dataframe.Sum},
    })
    agg = agg.Rename("Amount_sum", "Total")
    return agg
}

// pivotAggregation creates a simple pivot table from the aggregated dataframe.
func pivotAggregation(df dataframe.DataFrame, groupCol, pivotCol, valueCol string) dataframe.DataFrame {
    // Get unique group values
    groups := df.Col(groupCol).Unique()
    // Get unique pivot (month) values
    pivots := df.Col(pivotCol).Unique()

    // Build header row
    records := [][]string{{groupCol}}
    for i := 0; i < pivots.Len(); i++ {
        records[0] = append(records[0], fmt.Sprintf("Month_%v", pivots.Elem(i)))
    }

    // Build a lookup map: map[group][month] = total
    dataMap := make(map[string]map[string]string)
    for i := 0; i < df.Nrow(); i++ {
        row := df.Subset(i)
        groupVal := row.Col(groupCol).Elem(0)
        pivotVal := row.Col(pivotCol).Elem(0)
        totalVal := row.Col(valueCol).Elem(0)
        if _, ok := dataMap[groupVal]; !ok {
            dataMap[groupVal] = make(map[string]string)
        }
        dataMap[groupVal][pivotVal] = totalVal
    }

    // Build rows for each group
   	for i := 0; i < groups.Len(); i++ {
        groupVal := groups.Elem(i)
        row := []string{groupVal}
        for j := 0; j < pivots.Len(); j++ {
            pivotVal := pivots.Elem(j)
            if val, ok := dataMap[groupVal][pivotVal]; ok {
                row = append(row, val)
            } else {
                row = append(row, "0")
            }
        }
        records = append(records, row)
    }

    pivotDF := dataframe.LoadRecords(records)
    return pivotDF
}

func writeDataFrameToCSV(df dataframe.DataFrame, path string) {
    file, err := os.Create(path)
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    err = df.WriteCSV(file)
    if err != nil {
        fmt.Println("Error writing CSV:", err)
    }
}


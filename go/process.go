package main

import (
	"fmt"
	"os"
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
	dates := df.Col("Date").Records() // skip header row if present
	months := make([]int, len(dates))
	for i, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			months[i] = 0
		} else {
			months[i] = int(t.Month())
		}
	}
	monthSeries := series.New(months, series.Int, "Month")

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
	groups := []string{groupCol, monthCol}
	agg := df.GroupBy(groups...).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM},
		[]string{"Amount"},
	)
	agg = agg.Rename("Total", "Amount_SUM")
	return agg
}

func aggregateTwo(df dataframe.DataFrame, groupCol1, groupCol2, monthCol string) dataframe.DataFrame {
	groups := []string{groupCol1, groupCol2, monthCol}
	agg := df.GroupBy(groups...).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM},
		[]string{"Amount"},
	)
	agg = agg.Rename("Total", "Amount_SUM")
	return agg
}

// pivotAggregation creates a simple pivot table from the aggregated dataframe.
func pivotAggregation(df dataframe.DataFrame, groupCol, pivotCol, valueCol string) dataframe.DataFrame {
	// Get unique group values
	groups := make([]string, 0)
	groupColSeries := df.Col(groupCol)
	for i := 0; i < df.Nrow(); i++ {
		val := groupColSeries.Elem(i).String()
		if !contains(groups, val) {
			groups = append(groups, val)
		}
	}

	// Get unique pivot (month) values
	pivots := make([]string, 0)
	pivotColSeries := df.Col(pivotCol)
	for i := 0; i < df.Nrow(); i++ {
		val := pivotColSeries.Elem(i).String()
		if !contains(pivots, val) {
			pivots = append(pivots, val)
		}
	}

	// Build header row
	records := [][]string{{groupCol}}
	for _, pivot := range pivots {
		records[0] = append(records[0], fmt.Sprintf("Month_%v", pivot))
	}

	// Build a lookup map: map[string]map[string]string
	dataMap := make(map[string]map[string]string)
	for i := 0; i < df.Nrow(); i++ {
		row := df.Subset(i)
		groupVal := row.Col(groupCol).Elem(0).String()
		pivotVal := row.Col(pivotCol).Elem(0).String()
		totalVal := row.Col(valueCol).Elem(0).String()
		if _, ok := dataMap[groupVal]; !ok {
			dataMap[groupVal] = make(map[string]string)
		}
		dataMap[groupVal][pivotVal] = totalVal
	}

	// Build rows for each group
	for _, group := range groups {
		row := []string{group}
		for _, pivot := range pivots {
			if val, ok := dataMap[group][pivot]; ok {
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

// Helper function to check if a string slice contains a value
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func ProcessAccountMonthAggregation(csvPath string) bool {
	df := readAndPreprocessCSV(csvPath)
	if df.Err != nil {
		return false
	}

	accountAgg := aggregate(df, "Account", "Month")
	return writeDataFrameToCSV(accountAgg, "../results/go_account_month.csv")
}

func ProcessDepartmentMonthAggregation(csvPath string) bool {
	df := readAndPreprocessCSV(csvPath)
	if df.Err != nil {
		return false
	}

	departmentAgg := aggregate(df, "Department", "Month")
	return writeDataFrameToCSV(departmentAgg, "../results/go_department_month.csv")
}

func ProcessAccountDepartmentMonthAggregation(csvPath string) bool {
	df := readAndPreprocessCSV(csvPath)
	if df.Err != nil {
		return false
	}

	accountDeptAgg := aggregateTwo(df, "Account", "Department", "Month")
	return writeDataFrameToCSV(accountDeptAgg, "../results/go_account_dept_month.csv")
}

func CreatePivotTable(csvPath string) bool {
	df := readAndPreprocessCSV(csvPath)
	if df.Err != nil {
		return false
	}

	accountAgg := aggregate(df, "Account", "Month")
	pivot := pivotAggregation(accountAgg, "Account", "Month", "Total")
	return writeDataFrameToCSV(pivot, "../results/go_pivot_aggregation.csv")
}

func readAndPreprocessCSV(csvPath string) dataframe.DataFrame {
	// Open CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Error opening CSV:", err)
		return dataframe.DataFrame{Err: err}
	}
	defer file.Close()

	df := dataframe.ReadCSV(file)
	if df.Err != nil {
		fmt.Println("Error reading CSV:", df.Err)
		return df
	}

	// Add a Month column extracted from the Date column
	dates := df.Col("Date").Records()
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

	return df
}

func aggregate(df dataframe.DataFrame, groupCol string, monthCol string) dataframe.DataFrame {
	if df.Err != nil {
		return df
	}

	groups := []string{groupCol, monthCol}
	agg := df.GroupBy(groups...).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM},
		[]string{"Amount"},
	)
	agg = agg.Rename("Total", "Amount_SUM")
	return agg
}

func aggregateTwo(df dataframe.DataFrame, groupCol1, groupCol2, monthCol string) dataframe.DataFrame {
	if df.Err != nil {
		return df
	}

	groups := []string{groupCol1, groupCol2, monthCol}
	agg := df.GroupBy(groups...).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM},
		[]string{"Amount"},
	)
	agg = agg.Rename("Total", "Amount_SUM")
	return agg
}

func pivotAggregation(df dataframe.DataFrame, groupCol, pivotCol, valueCol string) dataframe.DataFrame {
	if df.Err != nil {
		return df
	}

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

	// Build data rows
	valueMap := make(map[string]map[string]string)
	for i := 0; i < df.Nrow(); i++ {
		group := groupColSeries.Elem(i).String()
		pivot := pivotColSeries.Elem(i).String()
		value := df.Col(valueCol).Elem(i).String()

		if valueMap[group] == nil {
			valueMap[group] = make(map[string]string)
		}
		valueMap[group][pivot] = value
	}

	for _, group := range groups {
		row := []string{group}
		for _, pivot := range pivots {
			value, exists := valueMap[group][pivot]
			if !exists {
				value = "0"
			}
			row = append(row, value)
		}
		records = append(records, row)
	}

	return dataframe.LoadRecords(records)
}

func writeDataFrameToCSV(df dataframe.DataFrame, path string) bool {
	if df.Err != nil {
		fmt.Printf("Error writing CSV: %v\n", df.Err)
		return false
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return false
	}
	defer f.Close()

	err = df.WriteCSV(f)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", path, err)
		return false
	}

	return true
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

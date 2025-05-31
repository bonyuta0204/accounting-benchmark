package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type DataFrameRunner struct {
	csvPath string
	df      dataframe.DataFrame
}

func (r *DataFrameRunner) Name() string {
	return "DataFrame"
}

func NewDataFrameRunner(csvPath string) *DataFrameRunner {
	return &DataFrameRunner{csvPath: csvPath}
}

func (r *DataFrameRunner) Run() error {
	// Load CSV
	start := time.Now()
	df, err := r.LoadCSV()
	if err != nil {
		return err
	}
	r.df = df

	fmt.Printf("CSV loading took: %s\n", time.Since(start))

	start = time.Now()
	r.AddMonthColumn()

	// Run aggregations
	// Account × Department × Monthly Aggregation
	agg := r.aggregateTwo("Account", "Department", "Month")
	if agg.Err != nil {
		return agg.Err
	}

	fmt.Printf("Aggregation took: %s\n", time.Since(start))

	start = time.Now()
	// Write to CSV
	if !r.writeDataFrameToCSV(agg, "../results/go_account_dept_month.csv") {
		return fmt.Errorf("failed to write CSV")
	}
	fmt.Printf("Write to CSV took: %s\n", time.Since(start))

	return nil
}

func (r *DataFrameRunner) LoadCSV() (dataframe.DataFrame, error) {
	file, err := os.Open(r.csvPath)
	if err != nil {
		return dataframe.DataFrame{}, err
	}
	defer file.Close()

	df := dataframe.ReadCSV(file)
	return df, df.Err
}

// AddMonthColumn adds a Month column to the DataFrame
func (r *DataFrameRunner) AddMonthColumn() {
	if r.df.Err != nil {
		return
	}

	// Add a Month column extracted from the Date column
	dates := r.df.Col("Date").Records()
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
	r.df = r.df.Mutate(monthSeries)
}

func (r *DataFrameRunner) aggregateTwo(groupCol1, groupCol2, monthCol string) dataframe.DataFrame {
	if r.df.Err != nil {
		return r.df
	}

	groups := []string{groupCol1, groupCol2, monthCol}
	agg := r.df.GroupBy(groups...).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM},
		[]string{"Amount"},
	)
	agg = agg.Rename("Total", "Amount_SUM")
	return agg
}

func (r *DataFrameRunner) writeDataFrameToCSV(df dataframe.DataFrame, path string) bool {
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

	err = r.df.WriteCSV(f)
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

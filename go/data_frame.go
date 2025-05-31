package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

var _ Aggregator = (*DataFrameRunner)(nil)

type DataFrameRunner struct {
	rawData    *dataframe.DataFrame
	resultData *dataframe.DataFrame
}

func NewDataFrameRunner() *DataFrameRunner {
	return &DataFrameRunner{}
}

func (r *DataFrameRunner) Name() string {
	return "DataFrame"
}

func (r *DataFrameRunner) LoadCSV(csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	rawData := dataframe.ReadCSV(file)
	r.rawData = &rawData
	return rawData.Err
}

func (r *DataFrameRunner) Aggregate(columns ...string) error {
	if r.rawData == nil {
		return fmt.Errorf("empty")

	}
	if r.rawData.Err != nil {
		return r.rawData.Err
	}

	// Add a Month column
	r.addMonthColumn()

	// Group by Date and Month
	result := r.rawData.GroupBy(columns...).Aggregation([]dataframe.AggregationType{
		dataframe.Aggregation_SUM,
	}, []string{"Amount"})

	r.resultData = &result

	return nil
}

// addMonthColumn adds a Month column to the DataFrame
func (r *DataFrameRunner) addMonthColumn() {
	if r.rawData.Err != nil {
		return
	}

	// Add a Month column extracted from the Date column
	dates := r.rawData.Col("Date").Records()
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
	rawData := r.rawData.Mutate(monthSeries)
	r.rawData = &rawData
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (r *DataFrameRunner) WriteToCSV(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if r.resultData == nil {
		return fmt.Errorf("empty")
	}
	if r.resultData.Err != nil {
		return r.resultData.Err
	}
	return r.resultData.WriteCSV(file)
}

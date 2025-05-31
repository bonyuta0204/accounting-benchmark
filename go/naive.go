package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

var _ Aggregator = (*NaiveAggregator)(nil)

type NaiveAggregator struct {
	csvPath    string
	headers    []string
	data       [][]string
	resultData [][]string
}

func NewNaiveAggregator() *NaiveAggregator {
	return &NaiveAggregator{}
}

func (r *NaiveAggregator) Name() string {
	return "Naive"
}

func (r *NaiveAggregator) LoadCSV(filepath string) error {
	r.csvPath = filepath
	
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}
	r.headers = headers
	
	// Read all data
	data, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	r.data = data
	
	// Add Month column by extracting from Date
	r.headers = append(r.headers, "Month")
	for i := range r.data {
		// Date format is YYYY-MM-DD, extract YYYY-MM as month
		if len(r.data[i][0]) >= 7 {
			r.data[i] = append(r.data[i], r.data[i][0][:7])
		} else {
			r.data[i] = append(r.data[i], "")
		}
	}
	
	return nil
}

func (r *NaiveAggregator) Aggregate(columns ...string) error {
	// Find column indices for groupBy columns and Amount
	groupByIndices := make([]int, len(columns))
	amountIndex := -1
	
	for i, col := range columns {
		found := false
		for j, header := range r.headers {
			if header == col {
				groupByIndices[i] = j
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("column %s not found", col)
		}
	}
	
	for j, header := range r.headers {
		if header == "Amount" {
			amountIndex = j
			break
		}
	}
	
	if amountIndex == -1 {
		return fmt.Errorf("Amount column not found")
	}
	
	// Aggregate using map
	aggregated := make(map[string]float64)
	
	for _, row := range r.data {
		// Build key from groupBy columns
		keyParts := make([]string, len(groupByIndices))
		for i, idx := range groupByIndices {
			keyParts[i] = row[idx]
		}
		key := strings.Join(keyParts, "|")
		
		// Parse amount
		amount, err := strconv.ParseFloat(row[amountIndex], 64)
		if err != nil {
			continue // Skip invalid amounts
		}
		
		aggregated[key] += amount
	}
	
	// Convert map to result slice
	result := make([][]string, 0, len(aggregated)+1)
	
	// Add headers
	headers := append(columns, "Amount_sum")
	result = append(result, headers)
	
	// Add data rows
	for key, sum := range aggregated {
		parts := strings.Split(key, "|")
		row := append(parts, fmt.Sprintf("%.2f", sum))
		result = append(result, row)
	}
	
	r.resultData = result
	return nil
}

func (r *NaiveAggregator) LoadParquet(filepath string) error {
	fr, err := local.NewLocalFileReader(filepath)
	if err != nil {
		return fmt.Errorf("failed to open parquet file: %w", err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(Transaction), 4)
	if err != nil {
		return fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer pr.ReadStop()

	num := int(pr.GetNumRows())
	transactions := make([]Transaction, num)
	if err = pr.Read(&transactions); err != nil {
		return fmt.Errorf("failed to read parquet data: %w", err)
	}

	// Convert to the same format as CSV data
	r.headers = []string{"Date", "Amount", "Account", "Department", "Month"}
	r.data = make([][]string, num)
	
	for i, t := range transactions {
		month := ""
		if len(t.Date) >= 7 {
			month = t.Date[:7]
		}
		r.data[i] = []string{
			t.Date,
			fmt.Sprintf("%.2f", t.Amount),
			t.Account,
			t.Department,
			month,
		}
	}
	
	return nil
}

func (r *NaiveAggregator) WriteToCSV(filepath string) error {
	if r.resultData == nil {
		return fmt.Errorf("no result data to write")
	}
	
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	for _, row := range r.resultData {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	
	return nil
}
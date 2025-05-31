package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
)

var _ Aggregator = (*DuckDBAggregator)(nil)

type DuckDBAggregator struct {
	db         *sql.DB
	csvPath    string
	resultRows [][]string
}

func NewDuckDBAggregator() *DuckDBAggregator {
	return &DuckDBAggregator{}
}

func (d *DuckDBAggregator) Name() string {
	return "DuckDB"
}

func (d *DuckDBAggregator) LoadCSV(filepath string) error {
	d.csvPath = filepath

	// Initialize DuckDB
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return fmt.Errorf("failed to open duckdb: %w", err)
	}
	d.db = db

	// Create table and load CSV
	createTableSQL := `
		CREATE TABLE transactions AS 
		SELECT *, 
		       strftime(Date, '%Y-%m') AS Month
		FROM read_csv_auto($1, header=true)
	`
	
	_, err = d.db.Exec(createTableSQL, filepath)
	if err != nil {
		return fmt.Errorf("failed to load CSV: %w", err)
	}

	return nil
}

func (d *DuckDBAggregator) Aggregate(columns ...string) error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Build the GROUP BY query
	groupByColumns := strings.Join(columns, ", ")
	selectColumns := groupByColumns + ", SUM(Amount) as Amount_sum"
	
	query := fmt.Sprintf(`
		SELECT %s
		FROM transactions
		GROUP BY %s
		ORDER BY %s
	`, selectColumns, groupByColumns, groupByColumns)

	rows, err := d.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to execute aggregation query: %w", err)
	}
	defer rows.Close()

	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// Store results
	d.resultRows = [][]string{cols}

	// Create a slice to hold the values
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Fetch all rows
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert to string slice
		row := make([]string, len(cols))
		for i, v := range values {
			if v == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		d.resultRows = append(d.resultRows, row)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	return nil
}

func (d *DuckDBAggregator) WriteToCSV(filepath string) error {
	if d.resultRows == nil || len(d.resultRows) == 0 {
		return fmt.Errorf("no result data to write")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range d.resultRows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// Close the database connection
func (d *DuckDBAggregator) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
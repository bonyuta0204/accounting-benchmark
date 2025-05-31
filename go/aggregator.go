package main

type Aggregator interface {
	// LoadCSV loads data from the given path.
	LoadCSV(path string) error
	// LoadParquet loads data from the given parquet file path.
	LoadParquet(path string) error
	// Aggregate runs aggregation on the loaded data.
	Aggregate(columns ...string) error
	// WriteToCSV writes the aggregated data to the given path.
	WriteToCSV(path string) error
	// Name returns the name of the aggregator.
	Name() string
}

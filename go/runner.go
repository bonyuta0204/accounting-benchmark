package main

type AggregationRunner interface {
	Run() error
	Name() string
}

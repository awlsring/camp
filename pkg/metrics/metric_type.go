package metrics

type Class int

const (
	Counter Class = iota
	Gauge
	Histogram
	Summary
)

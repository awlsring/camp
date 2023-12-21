package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Writer interface {
	PutMetric(metric string, class Class, value float64)
	Collect(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
}

type MetricWriterOpts func(*MetricWriter)

func WithNamespace(namespace string) MetricWriterOpts {
	return func(m *MetricWriter) {
		m.namespace = namespace
	}
}

type MetricWriter struct {
	namespace string
	counters  map[string]prometheus.Counter
	gauges    map[string]prometheus.Gauge
	histogram map[string]prometheus.Histogram
	summary   map[string]prometheus.Summary
}

func NewMetricWriter(opts ...MetricWriterOpts) Writer {
	m := &MetricWriter{
		namespace: "",
		counters:  make(map[string]prometheus.Counter),
		gauges:    make(map[string]prometheus.Gauge),
		histogram: make(map[string]prometheus.Histogram),
		summary:   make(map[string]prometheus.Summary),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *MetricWriter) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range m.counters {
		metric.Collect(ch)
	}
	for _, metric := range m.gauges {
		metric.Collect(ch)
	}
	for _, metric := range m.histogram {
		metric.Collect(ch)
	}
	for _, metric := range m.summary {
		metric.Collect(ch)
	}
}

func (m *MetricWriter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range m.counters {
		metric.Describe(ch)
	}
	for _, metric := range m.gauges {
		metric.Describe(ch)
	}
	for _, metric := range m.histogram {
		metric.Describe(ch)
	}
	for _, metric := range m.summary {
		metric.Describe(ch)
	}
}

func (m *MetricWriter) PutMetric(metric string, class Class, value float64) {
	switch class {
	case Counter:
		if _, ok := m.counters[metric]; !ok {
			m.counters[metric] = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: m.namespace,
				Name:      metric,
			})
		}
		m.counters[metric].Add(value)
	case Gauge:
		if _, ok := m.gauges[metric]; !ok {
			m.gauges[metric] = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: m.namespace,
				Name:      metric,
			})
		}
		m.gauges[metric].Set(value)
	case Histogram:
		if _, ok := m.histogram[metric]; !ok {
			m.histogram[metric] = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace: m.namespace,
				Name:      metric,
			})
		}
		m.histogram[metric].Observe(value)
	case Summary:
		if _, ok := m.summary[metric]; !ok {
			m.summary[metric] = prometheus.NewSummary(prometheus.SummaryOpts{
				Namespace: m.namespace,
				Name:      metric,
			})
		}
		m.summary[metric].Observe(value)
	}
}

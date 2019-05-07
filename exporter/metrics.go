package exporter

import "github.com/prometheus/client_golang/prometheus"

type Metrics map[string]*Metric

type Metric struct {
	Name      string
	Labels    []string
	ValueName string
	*prometheus.GaugeVec
}

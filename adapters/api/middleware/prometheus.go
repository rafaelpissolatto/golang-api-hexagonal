package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	prometheusClient "github.com/prometheus/client_model/go"
)

const (
	code                = "code"
	method              = "method"
	path                = "path"
	metricReqsTotal     = "requests_total"
	metricLatency       = "requests_duration_seconds"
	httpServerSubsystem = "http_server"
)

// CustomMetricRegistry custom prometheus metric registry
type CustomMetricRegistry struct {
	*prometheus.Registry
	customLabels []*prometheusClient.LabelPair
}

// Gather override the Gather function from prometheus to have custom labels
func (c *CustomMetricRegistry) Gather() ([]*prometheusClient.MetricFamily, error) {
	metricFamilies, err := c.Registry.Gather()
	for _, metricFamily := range metricFamilies {
		metrics := metricFamily.Metric
		for _, metric := range metrics {
			metric.Label = append(metric.Label, c.customLabels...)
		}
	}
	return metricFamilies, err
}

// NewCustomMetricsRegistry create a new custom prometheus metric registry
func NewCustomMetricsRegistry(labels map[string]string) *CustomMetricRegistry {
	custom := &CustomMetricRegistry{
		Registry: prometheus.NewRegistry(),
	}
	for key, value := range labels {
		k := key
		v := value
		custom.customLabels = append(custom.customLabels, &prometheusClient.LabelPair{
			Name:  &k,
			Value: &v,
		})
	}
	return custom
}

// NewPrometheusMiddleware create a new prometheus middleware metrics
func NewPrometheusMiddleware(serviceName string) *CustomMetricRegistry {
	// Set a label to all prometheus metrics
	customMetricRegistry := NewCustomMetricsRegistry(map[string]string{"service_name": serviceName})
	// Register default golang metrics on prometheus
	customMetricRegistry.MustRegister(collectors.NewGoCollector())
	return customMetricRegistry
}

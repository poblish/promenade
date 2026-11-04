package api

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	clientmodel "github.com/prometheus/client_model/go"
)

type TestHelper struct {
	metrics *PrometheusMetricsImpl
}

func (p *PrometheusMetricsImpl) TestHelper() *TestHelper {
	return &TestHelper{metrics: p}
}

func (h *TestHelper) Clear() {
	if h.metrics == nil {
		return
	}
	h.metrics.registry = prometheus.NewRegistry()
	h.metrics.registrations = newMetricRegistrations()
	h.metrics.errorCounterLock.Lock()
	h.metrics.errorCounter = nil
	h.metrics.errorCounterLock.Unlock()
}

func (h *TestHelper) Gather() ([]*clientmodel.MetricFamily, error) {
	if h.metrics == nil {
		return nil, nil
	}
	return h.metrics.registry.(prometheus.Gatherer).Gather()
}

func (h *TestHelper) MetricNames() []string {
	gathered, _ := h.Gather()

	count := len(gathered)
	names := make([]string, count)
	i := 0
	for _, entry := range gathered {
		names[i] = entry.GetName()
		i++
	}
	return names
}

type LabelsMap map[string]map[string]*clientmodel.Metric

func (h *TestHelper) GetMetricLabelValues(name string) LabelsMap {
	metricFamilies, _ := h.Gather()
	metricFamily, err := h.GetMetricFamily(metricFamilies, name)
	if err != nil {
		return nil
	}

	metrics := metricFamily.GetMetric()
	if err != nil {
		return nil
	}

	labelsMap := map[string]map[string]*clientmodel.Metric{}

	for _, m := range metrics {
		for _, label := range m.GetLabel() {
			labelMap := labelsMap[label.GetName()]
			if labelMap == nil {
				labelsMap[label.GetName()] = map[string]*clientmodel.Metric{label.GetValue(): m}
			} else {
				labelMap[label.GetValue()] = m
			}
		}
	}

	return labelsMap
}

func (h *TestHelper) GetMetricFamily(metrics []*clientmodel.MetricFamily, name string) (*clientmodel.MetricFamily, error) {
	for _, m := range metrics {
		if *m.Name == name {
			return m, nil
		}
	}

	return &clientmodel.MetricFamily{}, fmt.Errorf("metric %s not found in %v", name, metrics)
}

package api

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	clientmodel "github.com/prometheus/client_model/go"
)

type TestHelper struct {
	metrics *PrometheusMetrics
}

func (p *PrometheusMetrics) TestHelper() *TestHelper {
	return &TestHelper{metrics: p}
}

func (helper *TestHelper) Clear() {
	helper.metrics.registry = prometheus.NewRegistry()
	helper.metrics.registrations = newMetricRegistrations()
	helper.metrics.errorCounter = nil
}

func (helper *TestHelper) Gather() ([]*clientmodel.MetricFamily, error) {
	return helper.metrics.registry.(prometheus.Gatherer).Gather()
}

func (helper *TestHelper) MetricNames() []string {
	gathered, _ := helper.Gather()

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

func (helper *TestHelper) GetMetricLabelValues(name string) LabelsMap {
	metricFamilies, _ := helper.Gather()
	metricFamily, err := helper.GetMetricFamily(metricFamilies, name)
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

func (helper *TestHelper) GetMetricFamily(metrics []*clientmodel.MetricFamily, name string) (*clientmodel.MetricFamily, error) {
	for _, m := range metrics {
		if *m.Name == name {
			return m, nil
		}
	}

	return &clientmodel.MetricFamily{}, fmt.Errorf("metric %s not found in %v", name, metrics)
}

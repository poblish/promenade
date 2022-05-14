package api

import "github.com/prometheus/client_golang/prometheus"

type SummaryFacade struct {
	promMetric prometheus.Summary
}

var (
	DefaultObjectives = map[float64]float64{0.5: 0.01, 0.75: 0.01, 0.9: 0.01, 0.95: 0.01, 0.99: 0.01, 0.999: 0.01}
)

func (p *PrometheusMetricsImpl) buildSummary(builder MetricBuilder, name string, optionalDesc []string) SummaryFacade {
	return p.getOrAdd(name, TypeSummary, builder, optionalDesc).(SummaryFacade)
}

func (p *PrometheusMetricsImpl) Summary(name string, optionalDesc ...string) SummaryFacade {
	return p.buildSummary(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewSummary(prometheus.SummaryOpts{Name: fullMetricName, Help: fullDescription, Objectives: DefaultObjectives})
		p.Register(internal)
		return SummaryFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f SummaryFacade) Observe(value float64) {
	f.promMetric.Observe(value)
}

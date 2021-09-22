package api

import "github.com/prometheus/client_golang/prometheus"

type LabelledSummaryFacade struct {
	promMetric *prometheus.SummaryVec
}

func (p *PrometheusMetricsImpl) buildLabelledSummary(builder MetricBuilder, name string, optionalDesc []string) LabelledSummaryFacade {
	return p.getOrAdd(name, TypeSummaryLabels, builder, optionalDesc).(LabelledSummaryFacade)
}

func (p *PrometheusMetricsImpl) SummaryWithLabel(name string, labelName string, optionalDesc ...string) LabelledSummaryFacade {
	return p.SummaryWithLabels(name, []string{labelName}, optionalDesc...)
}

func (p *PrometheusMetricsImpl) SummaryWithLabels(name string, labelNames []string, optionalDesc ...string) LabelledSummaryFacade {
	return p.buildLabelledSummary(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: fullMetricName, Help: fullDescription, Objectives: DefaultObjectives}, labelNames)
		p.RegisterMetric(internal)
		return LabelledSummaryFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f LabelledSummaryFacade) Observe(value float64, labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Observe(value)
}

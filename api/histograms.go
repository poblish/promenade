package api

import "github.com/prometheus/client_golang/prometheus"

type HistogramFacade struct {
	promMetric prometheus.Histogram
}

var DefaultBuckets = prometheus.DefBuckets

func (p *PrometheusMetrics) buildHistogram(builder MetricBuilder, name string, optionalDesc []string) HistogramFacade {
	return p.getOrAdd(name, TypeHistogram, builder, optionalDesc).(HistogramFacade)
}

func (p *PrometheusMetrics) Histogram(name string, buckets []float64, optionalDesc ...string) HistogramFacade {
	return p.buildHistogram(func(p *PrometheusMetrics, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewHistogram(prometheus.HistogramOpts{Name: fullMetricName, Help: fullDescription, Buckets: buckets})
		p.RegisterMetric(internal)
		return HistogramFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (p *PrometheusMetrics) HistogramForResponseTime(name string, optionalDesc ...string) HistogramFacade {
	return p.buildHistogram(func(p *PrometheusMetrics, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewHistogram(prometheus.HistogramOpts{Name: fullMetricName, Help: fullDescription, Buckets: DefaultBuckets})
		p.RegisterMetric(internal)
		return HistogramFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f HistogramFacade) Update(value float64) {
	f.promMetric.Observe(value)
}

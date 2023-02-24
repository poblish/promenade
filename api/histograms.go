package api

import "github.com/prometheus/client_golang/prometheus"

type HistogramFacade interface {
	Update(value float64)
}

type histogramFacade struct {
	promMetric prometheus.Histogram
}

var DefaultBuckets = prometheus.DefBuckets

func (p *PrometheusMetricsImpl) buildHistogram(builder MetricBuilder, name string, optionalDesc []string) HistogramFacade {
	return p.getOrAdd(name, TypeHistogram, builder, optionalDesc).(HistogramFacade)
}

func (p *PrometheusMetricsImpl) Histogram(name string, buckets []float64, optionalDesc ...string) HistogramFacade {
	return p.buildHistogram(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewHistogram(prometheus.HistogramOpts{Name: fullMetricName, Help: fullDescription, Buckets: buckets})
		_ = p.Register(internal)
		return histogramFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (p *PrometheusMetricsImpl) HistogramForResponseTime(name string, optionalDesc ...string) HistogramFacade {
	return p.buildHistogram(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewHistogram(prometheus.HistogramOpts{Name: fullMetricName, Help: fullDescription, Buckets: DefaultBuckets})
		_ = p.Register(internal)
		return histogramFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f histogramFacade) Update(value float64) {
	f.promMetric.Observe(value)
}

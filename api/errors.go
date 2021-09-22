package api

import "github.com/prometheus/client_golang/prometheus"

type ErrorCounter struct {
	promMetric *prometheus.CounterVec
}

func (p *PrometheusMetricsImpl) Error(name string) ErrorCounter {
	return p.incrementError(name)
}

func (p *PrometheusMetricsImpl) incrementError(name string) ErrorCounter {
	var counter = p.getErrorCounter()
	counter.With(prometheus.Labels{"error_type": name}).Inc()
	return ErrorCounter{promMetric: counter}
}

// FIXME @Synchronized
func (p *PrometheusMetricsImpl) getErrorCounter() *prometheus.CounterVec {
	if p.errorCounter == nil {
		var adjustedName = p.metricNamePrefix + "errors"
		var description = adjustedName

		internal := prometheus.NewCounterVec(prometheus.CounterOpts{Name: adjustedName, Help: description}, []string{"error_type"})
		p.RegisterMetric(internal)
		p.errorCounter = internal
		p.errorCounterName = adjustedName
	}
	return p.errorCounter
}

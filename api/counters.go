package api

import "github.com/prometheus/client_golang/prometheus"

type CounterFacade interface {
	Inc()
	IncBy(inc float64)
}

type counterFacade struct {
	promMetric prometheus.Counter
}

func (p *PrometheusMetricsImpl) buildCounter(builder MetricBuilder, name string, optionalDesc []string) CounterFacade {
	return p.getOrAdd(name, TypeCounter, builder, optionalDesc).(CounterFacade)
}

func (p *PrometheusMetricsImpl) Counter(name string, optionalDesc ...string) CounterFacade {
	return p.buildCounter(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewCounter(prometheus.CounterOpts{Name: fullMetricName, Help: fullDescription})
		_ = p.Register(internal)
		return counterFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f counterFacade) Inc() {
	f.promMetric.Inc()
}

func (f counterFacade) IncBy(inc float64) {
	f.promMetric.Add(inc)
}

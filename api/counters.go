package api

import "github.com/prometheus/client_golang/prometheus"

type CounterFacade struct {
	promMetric prometheus.Counter
}

func (p *PrometheusMetricsImpl) buildCounter(builder MetricBuilder, name string, optionalDesc []string) CounterFacade {
	return p.getOrAdd(name, TypeCounter, builder, optionalDesc).(CounterFacade)
}

func (p *PrometheusMetricsImpl) Counter(name string, optionalDesc ...string) CounterFacade {
	return p.buildCounter(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewCounter(prometheus.CounterOpts{Name: fullMetricName, Help: fullDescription})
		p.Register(internal)
		return CounterFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f CounterFacade) Inc() {
	f.promMetric.Inc()
}

func (f CounterFacade) IncBy(inc float64) {
	f.promMetric.Add(inc)
}

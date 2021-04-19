package api

import "github.com/prometheus/client_golang/prometheus"

type CounterFacade struct {
	promMetric prometheus.Counter
}

func (p *PrometheusMetrics) buildCounter(builder MetricBuilder, name string, optionalDesc []string) CounterFacade {
	return p.getOrAdd(name, TypeCounter, builder, optionalDesc).(CounterFacade)
}

func (p *PrometheusMetrics) Counter(name string, optionalDesc ...string) CounterFacade {
	return p.buildCounter(func(p *PrometheusMetrics, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewCounter(prometheus.CounterOpts{Name: fullMetricName, Help: fullDescription})
		p.RegisterMetric(internal)
		return CounterFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f CounterFacade) Inc() {
	f.promMetric.Inc()
}

func (f CounterFacade) IncBy(inc float64) {
	f.promMetric.Add(inc)
}

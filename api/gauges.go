package api

import "github.com/prometheus/client_golang/prometheus"

type GaugeFacade struct {
	promMetric prometheus.Gauge
}

func (p *PrometheusMetricsImpl) buildGauge(builder MetricBuilder, name string, optionalDesc []string) GaugeFacade {
	return p.getOrAdd(name, TypeGauge, builder, optionalDesc).(GaugeFacade)
}

func (p *PrometheusMetricsImpl) Gauge(name string, optionalDesc ...string) GaugeFacade {
	return p.buildGauge(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewGauge(prometheus.GaugeOpts{Name: fullMetricName, Help: fullDescription})
		p.Register(internal)
		return GaugeFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f GaugeFacade) SetValue(value float64) {
	f.promMetric.Set(value)
}

func (f GaugeFacade) Inc() {
	f.promMetric.Inc()
}

func (f GaugeFacade) IncBy(inc float64) {
	f.promMetric.Add(inc)
}

func (f GaugeFacade) Dec() {
	f.promMetric.Dec()
}

func (f GaugeFacade) DecBy(dec float64) {
	f.promMetric.Sub(dec)
}

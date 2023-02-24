package api

import "github.com/prometheus/client_golang/prometheus"

type GaugeFacade interface {
	SetValue(value float64)
	Inc()
	IncBy(inc float64)
	Dec()
	DecBy(dec float64)
}

type gaugeFacade struct {
	promMetric prometheus.Gauge
}

func (p *PrometheusMetricsImpl) buildGauge(builder MetricBuilder, name string, optionalDesc []string) GaugeFacade {
	return p.getOrAdd(name, TypeGauge, builder, optionalDesc).(GaugeFacade)
}

func (p *PrometheusMetricsImpl) Gauge(name string, optionalDesc ...string) GaugeFacade {
	return p.buildGauge(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewGauge(prometheus.GaugeOpts{Name: fullMetricName, Help: fullDescription})
		_ = p.Register(internal)
		return gaugeFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f gaugeFacade) SetValue(value float64) {
	f.promMetric.Set(value)
}

func (f gaugeFacade) Inc() {
	f.promMetric.Inc()
}

func (f gaugeFacade) IncBy(inc float64) {
	f.promMetric.Add(inc)
}

func (f gaugeFacade) Dec() {
	f.promMetric.Dec()
}

func (f gaugeFacade) DecBy(dec float64) {
	f.promMetric.Sub(dec)
}

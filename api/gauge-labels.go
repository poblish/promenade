package api

import "github.com/prometheus/client_golang/prometheus"

type LabelledGaugeFacade struct {
	promMetric *prometheus.GaugeVec
}

func (p *PrometheusMetricsImpl) buildLabelledGauge(builder MetricBuilder, name string, optionalDesc []string) LabelledGaugeFacade {
	return p.getOrAdd(name, TypeGaugeLabels, builder, optionalDesc).(LabelledGaugeFacade)
}

func (p *PrometheusMetricsImpl) GaugeWithLabel(name string, labelName string, optionalDesc ...string) LabelledGaugeFacade {
	return p.GaugeWithLabels(name, []string{labelName}, optionalDesc...)
}

func (p *PrometheusMetricsImpl) GaugeWithLabels(name string, labelNames []string, optionalDesc ...string) LabelledGaugeFacade {
	return p.buildLabelledGauge(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: fullMetricName, Help: fullDescription}, labelNames)
		p.Register(internal)
		return LabelledGaugeFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f LabelledGaugeFacade) IncLabels(labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Inc()
}

func (f LabelledGaugeFacade) DecLabels(labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Dec()
}

func (f LabelledGaugeFacade) IncLabelsBy(labelValues ...string) IncGaugeByValue {
	return IncGaugeByValue{gauge: f.promMetric.WithLabelValues(labelValues...)}
}

func (f LabelledGaugeFacade) DecLabelsBy(labelValues ...string) DecGaugeByValue {
	return DecGaugeByValue{gauge: f.promMetric.WithLabelValues(labelValues...)}
}

func (f LabelledGaugeFacade) SetLabels(labelValues ...string) SetGaugeByValue {
	return SetGaugeByValue{gauge: f.promMetric.WithLabelValues(labelValues...)}
}

type IncGaugeByValue struct {
	gauge prometheus.Gauge
}

type DecGaugeByValue struct {
	gauge prometheus.Gauge
}

type SetGaugeByValue struct {
	gauge prometheus.Gauge
}

func (f IncGaugeByValue) Value(inc float64) {
	f.gauge.Add(inc)
}

func (f DecGaugeByValue) Value(inc float64) {
	f.gauge.Sub(inc)
}

func (f SetGaugeByValue) Value(inc float64) {
	f.gauge.Set(inc)
}

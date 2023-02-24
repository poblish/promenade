package api

import "github.com/prometheus/client_golang/prometheus"

type LabelledGaugeFacade interface {
	IncLabels(labelValues ...string)
	IncLabelsBy(labelValues ...string) IncGaugeByValue
	DecLabels(labelValues ...string)
	DecLabelsBy(labelValues ...string) DecGaugeByValue
	SetLabels(labelValues ...string) SetGaugeByValue
}

type labelledGaugeFacade struct {
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
		_ = p.Register(internal)
		return labelledGaugeFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f labelledGaugeFacade) IncLabels(labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Inc()
}

func (f labelledGaugeFacade) DecLabels(labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Dec()
}

func (f labelledGaugeFacade) IncLabelsBy(labelValues ...string) IncGaugeByValue {
	return incGaugeByValue{gauge: f.promMetric.WithLabelValues(labelValues...)}
}

func (f labelledGaugeFacade) DecLabelsBy(labelValues ...string) DecGaugeByValue {
	return decGaugeByValue{gauge: f.promMetric.WithLabelValues(labelValues...)}
}

func (f labelledGaugeFacade) SetLabels(labelValues ...string) SetGaugeByValue {
	return setGaugeByValue{gauge: f.promMetric.WithLabelValues(labelValues...)}
}

type IncGaugeByValue interface {
	Value(inc float64)
}

type incGaugeByValue struct {
	gauge prometheus.Gauge
}

type DecGaugeByValue interface {
	Value(inc float64)
}

type decGaugeByValue struct {
	gauge prometheus.Gauge
}

type SetGaugeByValue interface {
	Value(inc float64)
}

type setGaugeByValue struct {
	gauge prometheus.Gauge
}

func (f incGaugeByValue) Value(inc float64) {
	f.gauge.Add(inc)
}

func (f decGaugeByValue) Value(inc float64) {
	f.gauge.Sub(inc)
}

func (f setGaugeByValue) Value(inc float64) {
	f.gauge.Set(inc)
}

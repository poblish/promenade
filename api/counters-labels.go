package api

import "github.com/prometheus/client_golang/prometheus"

type LabelledCounterFacade interface {
	IncLabel(labelValues ...string)
	IncLabelBy(labelValues ...string) IncByValue
}

type labelledCounterFacade struct {
	promMetric *prometheus.CounterVec
}

func (p *PrometheusMetricsImpl) buildLabelledCounter(builder MetricBuilder, name string, optionalDesc []string) LabelledCounterFacade {
	return p.getOrAdd(name, TypeCounterLabels, builder, optionalDesc).(LabelledCounterFacade)
}

func (p *PrometheusMetricsImpl) CounterWithLabel(name string, labelName string, optionalDesc ...string) LabelledCounterFacade {
	return p.CounterWithLabels(name, []string{labelName}, optionalDesc...)
}

func (p *PrometheusMetricsImpl) CounterWithLabels(name string, labelNames []string, optionalDesc ...string) LabelledCounterFacade {
	return p.buildLabelledCounter(func(p *PrometheusMetricsImpl, fullMetricName string, fullDescription string) interface{} {
		internal := prometheus.NewCounterVec(prometheus.CounterOpts{Name: fullMetricName, Help: fullDescription}, labelNames)
		_ = p.Register(internal)
		return labelledCounterFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f labelledCounterFacade) IncLabel(labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Inc()
}

type IncByValue interface {
	Value(inc float64)
}

type incByValue struct {
	counter prometheus.Counter
}

func (f labelledCounterFacade) IncLabelBy(labelValues ...string) IncByValue {
	return incByValue{counter: f.promMetric.WithLabelValues(labelValues...)}
}

func (f incByValue) Value(inc float64) {
	f.counter.Add(inc)
}

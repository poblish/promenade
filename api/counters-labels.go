package api

import "github.com/prometheus/client_golang/prometheus"

type LabelledCounterFacade struct {
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
		p.Register(internal)
		return LabelledCounterFacade{promMetric: internal}
	}, name, optionalDesc)
}

func (f LabelledCounterFacade) IncLabel(labelValues ...string) {
	f.promMetric.WithLabelValues(labelValues...).Inc()
}

type IncByValue struct {
	counter prometheus.Counter
}

func (f LabelledCounterFacade) IncLabelBy(labelValues ...string) IncByValue {
	return IncByValue{counter: f.promMetric.WithLabelValues(labelValues...)}
}

func (f IncByValue) Value(inc float64) {
	f.counter.Add(inc)
}

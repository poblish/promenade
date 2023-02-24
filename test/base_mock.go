package test

import (
	"github.com/poblish/promenade/api"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type BaseMock struct {
}

func (b BaseMock) Register(metric prometheus.Collector) error {
	log.Print("Registering", metric)
	return nil
}

func (b BaseMock) MustRegister(metric prometheus.Collector) {
	_ = b.Register(metric)
}

func (b BaseMock) TestHelper() *api.TestHelper {
	return &api.TestHelper{}
}

func (b BaseMock) Counter(name string, optionalDesc ...string) api.CounterFacade {
	return api.CounterFacade{}
}

func (b BaseMock) CounterWithLabel(name string, labelName string, optionalDesc ...string) api.LabelledCounterFacade {
	return api.LabelledCounterFacade{}
}

func (b BaseMock) CounterWithLabels(name string, labelNames []string, optionalDesc ...string) api.LabelledCounterFacade {
	return api.LabelledCounterFacade{}
}

func (b BaseMock) Error(name string) api.ErrorCounter {
	return api.ErrorCounter{}
}

func (b BaseMock) Gauge(name string, optionalDesc ...string) api.GaugeFacade {
	return api.GaugeFacade{}
}

func (b BaseMock) GaugeWithLabel(name string, labelName string, optionalDesc ...string) api.LabelledGaugeFacade {
	return api.LabelledGaugeFacade{}
}

func (b BaseMock) GaugeWithLabels(name string, labelNames []string, optionalDesc ...string) api.LabelledGaugeFacade {
	return api.LabelledGaugeFacade{}
}

func (b BaseMock) Histogram(name string, buckets []float64, optionalDesc ...string) api.HistogramFacade {
	return api.HistogramFacade{}
}

func (b BaseMock) HistogramForResponseTime(name string, optionalDesc ...string) api.HistogramFacade {
	return api.HistogramFacade{}
}

func (b BaseMock) Summary(name string, optionalDesc ...string) api.SummaryFacade {
	return api.SummaryFacade{}
}

func (b BaseMock) SummaryWithLabel(name string, labelName string, optionalDesc ...string) api.LabelledSummaryFacade {
	return api.LabelledSummaryFacade{}
}

func (b BaseMock) SummaryWithLabels(name string, labelNames []string, optionalDesc ...string) api.LabelledSummaryFacade {
	return api.LabelledSummaryFacade{}
}

func (b BaseMock) Timer(Name string) func() time.Duration {
	return emptyTimer
}

func (b BaseMock) TimerWithLabel(Name string, labelName string, labelValue string) func() time.Duration {
	return emptyTimer
}
func (b BaseMock) TimerWithLabels(Name string, labelNames []string, labelValues []string) func() time.Duration {
	return emptyTimer
}

func emptyTimer() time.Duration {
	return time.Duration(0)
}

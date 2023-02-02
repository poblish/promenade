package api

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestTimersControlled(t *testing.T) {
	metrics := PrometheusMetricsImpl{registry: prometheus.NewRegistry(),
		metricNamePrefix: "xx_",
		registrations:    newMetricRegistrations(),
		normalisedNames:  newNormalisedNames(),
		timerFactory:     &controlledTimerFactory{defaultExpectation: 2 * time.Second}}

	timedMethod(&metrics)
	timedMethod(&metrics)

	m := findMetric("xx_timer", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "summary:<sample_count:2 sample_sum:4 quantile:<quantile:0.5 value:2 > quantile:<quantile:0.75 value:2 > quantile:<quantile:0.9 value:2 > quantile:<quantile:0.95 value:2 > quantile:<quantile:0.99 value:2 > quantile:<quantile:0.999 value:2 > >",
		strings.TrimSpace(m.Metric[0].String()))
}

func TestTimersRealTime(t *testing.T) {
	metrics := PrometheusMetricsImpl{registry: prometheus.NewRegistry(),
		registrations:   newMetricRegistrations(),
		normalisedNames: newNormalisedNames(),
		timerFactory:    &defaultTimerFactory{}}

	timedMethod(&metrics)

	m := findMetric("timer", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	// Can't validate times
}

func TestLabelledTimersControlled(t *testing.T) {
	metrics := PrometheusMetricsImpl{registry: prometheus.NewRegistry(),
		metricNamePrefix: "xx_",
		registrations:    newMetricRegistrations(),
		normalisedNames:  newNormalisedNames(),
		timerFactory:     &controlledTimerFactory{defaultExpectation: 2 * time.Second}}

	timedMethodWithLabel(&metrics)
	timedMethodWithLabel(&metrics)

	m := findMetric("xx_animal_timer", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > summary:<sample_count:2 sample_sum:4 quantile:<quantile:0.5 value:2 > quantile:<quantile:0.75 value:2 > quantile:<quantile:0.9 value:2 > quantile:<quantile:0.95 value:2 > quantile:<quantile:0.99 value:2 > quantile:<quantile:0.999 value:2 > >",
		strings.TrimSpace(m.Metric[0].String()))
}

func TestMultiLabelledTimersControlled(t *testing.T) {
	metrics := PrometheusMetricsImpl{registry: prometheus.NewRegistry(),
		metricNamePrefix: "xx_",
		registrations:    newMetricRegistrations(),
		normalisedNames:  newNormalisedNames(),
		timerFactory:     &controlledTimerFactory{defaultExpectation: 2 * time.Second}}

	timedMethodWithLabels(&metrics)
	timedMethodWithLabels(&metrics)

	m := findMetric("xx_animal_timer", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > label:<name:\"colour\" value:\"yellow\" > label:<name:\"feet\" value:\"4\" > summary:<sample_count:2 sample_sum:4 quantile:<quantile:0.5 value:2 > quantile:<quantile:0.75 value:2 > quantile:<quantile:0.9 value:2 > quantile:<quantile:0.95 value:2 > quantile:<quantile:0.99 value:2 > quantile:<quantile:0.999 value:2 > >",
		strings.TrimSpace(m.Metric[0].String()))
}

func timedMethod(metrics PrometheusMetrics) {
	defer metrics.Timer("Timer")()
	fmt.Println("Whatever it is we're timing")
}

func timedMethodWithLabel(metrics PrometheusMetrics) {
	defer metrics.TimerWithLabel("animal_timer", "animal", "cat")()
	fmt.Println("Whatever it is we're timing")
}

func timedMethodWithLabels(metrics PrometheusMetrics) {
	defer metrics.TimerWithLabels("animal_timer", []string{"animal", "feet", "colour"}, []string{"cat", "4", "yellow"})()
	fmt.Println("Whatever it is we're timing")
}

type controlledTimerFactory struct {
	timerFactory
	defaultExpectation time.Duration
}

type controlledTimer struct {
	eventTimer
	observer    prometheus.Observer
	expectation time.Duration
}

func (f *controlledTimerFactory) NewTimer(o prometheus.Observer) eventTimer {
	return &controlledTimer{observer: o, expectation: f.defaultExpectation}
}

func (t *controlledTimer) Observe() time.Duration {
	fmt.Println("Making fixed observation...")
	t.observer.Observe(t.expectation.Seconds())
	return t.expectation
}

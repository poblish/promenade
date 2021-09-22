package api

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type eventTimer interface {
	Observe() time.Duration
}

type timerFactory interface {
	NewTimer(o prometheus.Observer) eventTimer
}

type defaultTimer struct {
	eventTimer
	timer *prometheus.Timer
}

type defaultTimerFactory struct {
	timerFactory
}

func (f *defaultTimerFactory) NewTimer(o prometheus.Observer) eventTimer {
	return &defaultTimer{timer: prometheus.NewTimer(o)}
}

func (t *defaultTimer) Observe() time.Duration {
	return t.timer.ObserveDuration()
}

func (p *PrometheusMetricsImpl) Timer(Name string) func() time.Duration {
	timer := p.timerFactory.NewTimer(p.Summary(Name).promMetric)
	return func() time.Duration {
		diff := timer.Observe()
		// fmt.Println("Observing", diff)
		return diff
	}
}

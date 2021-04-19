package api

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

var caseInsensitiveMetrics = NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah"})
var caseSensitiveMetrics = NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah", CaseSensitiveMetricNames: true})

func BenchmarkCounterCaseSensitive(b *testing.B) {
	for n := 0; n < b.N; n++ {
		caseSensitiveMetrics.Counter("ANonReuseS").Inc()
	}
}

func BenchmarkCounterCaseInsensitive(b *testing.B) {
	for n := 0; n < b.N; n++ {
		caseInsensitiveMetrics.Counter("ANonReuseI").Inc()
	}
}

func BenchmarkCounterReuse(b *testing.B) {
	x := caseInsensitiveMetrics.Counter("AReuse")
	for n := 0; n < b.N; n++ {
		x.Inc()
	}
}

func BenchmarkLabelledCounterReuse(b *testing.B) {
	x := caseInsensitiveMetrics.CounterWithLabels("Labelled", []string{"country"})
	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			x.IncLabel("uk")
		} else {
			x.IncLabel("usa")
		}
	}
}

func BenchmarkSummary(b *testing.B) {
	x := caseInsensitiveMetrics.Summary("MySummary")
	for n := 0; n < b.N; n++ {
		x.Observe(float64(n))
	}
}

func BenchmarkTimer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		anotherTimedMethod()
	}
}

func anotherTimedMethod() {
	defer caseInsensitiveMetrics.Timer("T")()
}

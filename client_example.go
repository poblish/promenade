package main

import (
	"fmt"

	promenade "github.com/poblish/promenade/api"
)

func main() {
	metrics := promenade.NewMetrics(promenade.MetricOpts{MetricNamePrefix: "prefix"})
	metrics.Counter("c")
	metrics.CounterWithLabel("places", "city").IncLabel("London")
	metrics.CounterWithLabels("animals", []string{"type", "breed"}).IncLabel("cat", "persian")
	metrics.Error("e")
	metrics.Gauge("g")
	metrics.HistogramForResponseTime("h")
	metrics.Histogram("hb", []float64{1, 10})
	metrics.Summary("s")
	timedMethod(&metrics)

	fmt.Println(metrics.TestHelper().MetricNames())
}

func timedMethod(metrics promenade.PrometheusMetrics) {
	defer metrics.Timer("t")()
	fmt.Println("Whatever it is we're timing")
}

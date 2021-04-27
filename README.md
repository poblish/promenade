# Promenade

A simplified, and slightly more opinionated Prometheus client API for Golang.

The goal is to encourage much greater metric use within a project without excessive lines of code, and with much of the configurability and flexibility of the [official Prometheus client](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) API tucked away.

Registration happens on first use. All metric names are normalised and (by default) reduced to lowercase.

----

## Examples

```golang

import (
	"fmt"

	promenade "github.com/poblish/promenade/api"
)

func main() {
    metrics := promenade.NewMetrics(promenade.MetricOpts{MetricNamePrefix: "prefix"})

    metrics.Counter("c").Inc() // Increment prefix_c

    // Increment {city:London} label for prefix_places
    metrics.CounterWithLabel("places", "city").IncLabel("London")

    // Increment {type:cat, breed:persian} labels for prefix_animals
    metrics.CounterWithLabels("animals", []string{"type", "breed"}).IncLabel("cat", "persian")

    // Increment {error_type:bade} label for prefix_errors
    metrics.Error("bad")

    // Histograms
    histograms(&metrics)
    histogram_buckets(&metrics)

    // Timers
    timedMethod(&metrics)
}

func histograms(metrics *promenade.PrometheusMetrics) {
	times := metrics.HistogramForResponseTime("latency")
    times.Update(0.03)
    times.Update(0.05)
}

func histogram_buckets(metrics *promenade.PrometheusMetrics) {
	ages := metrics.Histogram("population_by_age", []float64{18, 25, 25, 45, 55, 65})
    ages.Update(21)
    ages.Update(45)
    ages.Update(81)
}

func timedMethod(metrics *promenade.PrometheusMetrics) {
	defer metrics.Timer("calculate Pi")()  // Start the timer, observe on exit

	fmt.Println("Start doing it...")
    // ...
}

```
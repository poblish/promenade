package test

import (
	"fmt"
	"github.com/poblish/promenade/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMockApiSpec(t *testing.T) {
	test(t, &BaseMock{})
}

func test(t *testing.T, m api.PrometheusMetrics) {
	require.NotNil(t, m)

	m.Counter("c").Inc() // Increment prefix_c

	// Increment {city:London} label for prefix_places
	m.CounterWithLabel("places", "city").IncLabel("London")

	// Increment {type:cat, breed:persian} labels for prefix_animals
	m.CounterWithLabels("animals", []string{"type", "breed"}).IncLabel("cat", "persian")

	// Gauges
	m.Gauge("g").SetValue(101)
	m.Gauge("g").Dec()

	m.GaugeWithLabel("a", "b", "c").IncLabels("b")
	m.GaugeWithLabels("current animals", []string{"animal", "breed"}).SetLabels("fleas", "plague").Value(500)

	// Increment {error_type:bad} label for prefix_errors
	m.Error("bad")

	// Histograms
	histograms(m)
	histogram_buckets(m)

	// Timers
	timedMethod(m)

	// Tests
	testMethods(m)
}

func histograms(m api.PrometheusMetrics) {
	times := m.HistogramForResponseTime("latency")
	times.Update(0.03)
	times.Update(0.05)
}

func histogram_buckets(m api.PrometheusMetrics) {
	ages := m.Histogram("population_by_age", []float64{18, 25, 25, 45, 55, 65})
	ages.Update(21)
	ages.Update(45)
	ages.Update(81)
}

func timedMethod(m api.PrometheusMetrics) {
	defer m.Timer("calculate Pi")() // Start the timer, observe on exit

	fmt.Println("Start doing it...")
	// ...
}

func testMethods(m api.PrometheusMetrics) {
	m.TestHelper().Clear()  // reset; start with new registry
	m.TestHelper().Gather() // gather all registered Collectors
	m.TestHelper().MetricNames()
	// etc.
}

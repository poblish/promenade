package api

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	clientmodel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "z"})
	doTestCounter(t, &metrics)
}

func doTestCounter(t *testing.T, metrics *PrometheusMetricsImpl) {
	c := metrics.Counter("Mine")
	c.Inc()
	metrics.Counter("Mine").Inc()
	c.IncBy(7)

	m := findMetric("z_mine", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "name:\"z_mine\" help:\"z_mine\" type:COUNTER metric:<counter:<value:9 > >", strings.TrimSpace(m.String()))
}

func TestCounterWithExplicitDescription(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah"})

	metrics.Counter("123", "MyDesc").Inc()

	m := findMetric("blah_123", metrics.gatherOK(t))
	assert.Contains(t, m.String(), "name:\"blah_123\" help:\"MyDesc\" type:COUNTER")
}

func TestCounterWithMappedDescription(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(),
		MetricNamePrefix: "a",
		Descriptions:     MetricDescriptions{"mapped": "Description found", "mapped_blank": ""}})

	metrics.Counter("mapped").Inc()
	metrics.Counter("unmapped").Inc()
	metrics.Counter("mapped_blank").Inc()

	gathered := metrics.gatherOK(t)
	assert.Contains(t, findMetric("a_mapped", gathered).String(), "name:\"a_mapped\" help:\"Description found\" type:COUNTER")
	assert.Contains(t, findMetric("a_unmapped", gathered).String(), "name:\"a_unmapped\" help:\"a_unmapped\" type:COUNTER")
	assert.Contains(t, findMetric("a_mapped_blank", gathered).String(), "name:\"a_mapped_blank\" help:\"\" type:COUNTER")
}

func TestCounterCaseInsensitivity(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah"})

	metrics.Counter("abcd").Inc()
	metrics.Counter("AbCd").Inc()
	metrics.Counter("ABCD").Inc()

	m := findMetric("blah_abcd", metrics.gatherOK(t))
	assert.Equal(t, "name:\"blah_abcd\" help:\"blah_abcd\" type:COUNTER metric:<counter:<value:3 > >", strings.TrimSpace(m.String()))
}

func TestCounterCaseSensitivity(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah", CaseSensitiveMetricNames: true})

	metrics.Counter("abcd").Inc()
	metrics.Counter("AbCd").Inc()
	metrics.Counter("ABCD").Inc()

	gathered := metrics.gatherOK(t)
	assert.Equal(t, "name:\"blah_abcd\" help:\"blah_abcd\" type:COUNTER metric:<counter:<value:1 > >", strings.TrimSpace(findMetric("blah_abcd", gathered).String()))
	assert.Equal(t, "name:\"blah_AbCd\" help:\"blah_AbCd\" type:COUNTER metric:<counter:<value:1 > >", strings.TrimSpace(findMetric("blah_AbCd", gathered).String()))
	assert.Equal(t, "name:\"blah_ABCD\" help:\"blah_ABCD\" type:COUNTER metric:<counter:<value:1 > >", strings.TrimSpace(findMetric("blah_ABCD", gathered).String()))
}

func TestCounterWithAllDefaultOptions(t *testing.T) {
	metrics := NewMetrics(MetricOpts{})

	metrics.Counter("abcde").Inc()

	m := findMetric("abcde", metrics.gatherOK(t))
	assert.Equal(t, "name:\"abcde\" help:\"abcde\" type:COUNTER metric:<counter:<value:1 > >", strings.TrimSpace(m.String()))
}

func TestCounterWithBlankExplicitName(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah"})

	metrics.Counter("234", "").Inc()

	m := findMetric("blah_234", metrics.gatherOK(t))
	assert.Equal(t, "name:\"blah_234\" help:\"blah_234\" type:COUNTER metric:<counter:<value:1 > >", strings.TrimSpace(m.String()))
}

func TestCounterWithLabel(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "v"})

	c := metrics.CounterWithLabel("visitors", "country", "desc")
	c.IncLabel("uk")
	c.IncLabelBy("usa").Value(16)
	c.IncLabel("uk")
	c.IncLabel("usa")
	c.IncLabelBy("usa").Value(3)

	labels := metrics.TestHelper().GetMetricLabelValues("v_visitors")
	assertRegistryHasMetricWithConfigTypeCounts(t, labels, "v_visitors", map[string]map[string]float64{"country": {
		"uk":  2.0,
		"usa": 20.0,
	}})

	// For coverage
	assert.Nil(t, metrics.TestHelper().GetMetricLabelValues("_xxx"))
}
func TestCounterWithLabels(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "v", PrefixSeparator: ":"})

	c := metrics.CounterWithLabels("animals", []string{"animal", "breed"}, "desc")
	c.IncLabel("cat", "persian")
	c.IncLabelBy("dog", "spaniel").Value(16)
	c.IncLabel("cat", "black")
	c.IncLabel("dog", "greyhound")
	c.IncLabelBy("cat", "black").Value(3)

	m := findMetric("v:animals", metrics.gatherOK(t))
	assert.Equal(t, "name:\"v:animals\" help:\"desc\" type:COUNTER metric:<label:<name:\"animal\" value:\"cat\" > label:<name:\"breed\" value:\"black\" > counter:<value:4 > > metric:<label:<name:\"animal\" value:\"cat\" > label:<name:\"breed\" value:\"persian\" > counter:<value:1 > > metric:<label:<name:\"animal\" value:\"dog\" > label:<name:\"breed\" value:\"greyhound\" > counter:<value:1 > > metric:<label:<name:\"animal\" value:\"dog\" > label:<name:\"breed\" value:\"spaniel\" > counter:<value:16 > >", strings.TrimSpace(m.String()))
}

func TestBadNameReuse(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah"})
	metrics.Counter("a").Inc()
	assert.Panics(t, func() { metrics.Gauge("a").Inc() }, "The code did not panic")
}

func TestGauge(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "x-service#123"})

	g := metrics.Gauge("MyGauge")
	g.SetValue(100)
	g.Inc()
	g.IncBy(5)
	g.Dec()
	g.DecBy(4)
	metrics.Gauge("MyGauge").Inc()

	m := findMetric("x_service_123_mygauge", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "name:\"x_service_123_mygauge\" help:\"x_service_123_mygauge\" type:GAUGE metric:<gauge:<value:102 > >", strings.TrimSpace(m.String()))
}

func TestGaugeWithLabel(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "prefix"})

	g := metrics.GaugeWithLabel("current animals", "animal")
	g.SetLabels("fleas").Value(1000)
	g.IncLabels("dog")
	g.IncLabelsBy("cat").Value(5)
	g.DecLabels("cat")
	g.DecLabelsBy("fleas").Value(15)

	m := findMetric("prefix_current_animals", metrics.gatherOK(t))
	assert.Equal(t, 3, len(m.Metric))
	assert.Equal(t, "name:\"prefix_current_animals\" help:\"prefix_current_animals\" type:GAUGE metric:<label:<name:\"animal\" value:\"cat\" > gauge:<value:4 > > metric:<label:<name:\"animal\" value:\"dog\" > gauge:<value:1 > > metric:<label:<name:\"animal\" value:\"fleas\" > gauge:<value:985 > >", strings.TrimSpace(m.String()))
}

func TestGaugeWithLabels(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "prefix"})

	g := metrics.GaugeWithLabels("current animals", []string{"animal", "breed"})
	g.SetLabels("fleas", "plague").Value(500)
	g.SetLabels("fleas", "asian").Value(1000)
	g.IncLabels("dog", "borzoi")
	g.IncLabelsBy("cat", "black").Value(5)
	g.IncLabelsBy("cat", "white").Value(1)
	g.DecLabels("cat", "black")
	g.DecLabels("cat", "white")
	g.DecLabelsBy("fleas", "plague").Value(15)

	m := findMetric("prefix_current_animals", metrics.gatherOK(t))
	assert.Equal(t, 5, len(m.Metric))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > label:<name:\"breed\" value:\"black\" > gauge:<value:4 >",
		strings.TrimSpace(m.Metric[0].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > label:<name:\"breed\" value:\"white\" > gauge:<value:0 >",
		strings.TrimSpace(m.Metric[1].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"dog\" > label:<name:\"breed\" value:\"borzoi\" > gauge:<value:1 >",
		strings.TrimSpace(m.Metric[2].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"fleas\" > label:<name:\"breed\" value:\"asian\" > gauge:<value:1000 >",
		strings.TrimSpace(m.Metric[3].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"fleas\" > label:<name:\"breed\" value:\"plague\" > gauge:<value:485 >",
		strings.TrimSpace(m.Metric[4].String()))
}

func TestSummary(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "BLAH"})

	s := metrics.Summary("MySummary")
	s.Observe(1.3)
	s.Observe(2.5)
	s.Observe(2.6)
	s.Observe(2.9)
	s.Observe(3.2)
	s.Observe(3.3)
	s.Observe(3.834344)

	m := findMetric("blah_mysummary", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "summary:<sample_count:7 sample_sum:19.634344000000002 quantile:<quantile:0.5 value:2.9 > quantile:<quantile:0.75 value:3.3 > quantile:<quantile:0.9 value:3.834344 > quantile:<quantile:0.95 value:3.834344 > quantile:<quantile:0.99 value:3.834344 > quantile:<quantile:0.999 value:3.834344 > >",
		strings.TrimSpace(m.Metric[0].String()))
}

func TestSummaryWithLabel(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "BLAH"})

	s := metrics.SummaryWithLabel("animal facts", "animal")
	s.Observe(1.0, "cat")
	s.Observe(2.5, "cat")
	s.Observe(2.6, "dog")
	s.Observe(2.0, "cat")
	s.Observe(3.2, "ant")
	s.Observe(3.3, "dog")
	s.Observe(3.834344, "bear")

	m := findMetric("blah_animal_facts", metrics.gatherOK(t))
	assert.Equal(t, 4, len(m.Metric))
	assert.Equal(t, "label:<name:\"animal\" value:\"ant\" > summary:<sample_count:1 sample_sum:3.2 quantile:<quantile:0.5 value:3.2 > quantile:<quantile:0.75 value:3.2 > quantile:<quantile:0.9 value:3.2 > quantile:<quantile:0.95 value:3.2 > quantile:<quantile:0.99 value:3.2 > quantile:<quantile:0.999 value:3.2 > >",
		strings.TrimSpace(m.Metric[0].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"bear\" > summary:<sample_count:1 sample_sum:3.834344 quantile:<quantile:0.5 value:3.834344 > quantile:<quantile:0.75 value:3.834344 > quantile:<quantile:0.9 value:3.834344 > quantile:<quantile:0.95 value:3.834344 > quantile:<quantile:0.99 value:3.834344 > quantile:<quantile:0.999 value:3.834344 > >",
		strings.TrimSpace(m.Metric[1].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > summary:<sample_count:3 sample_sum:5.5 quantile:<quantile:0.5 value:2 > quantile:<quantile:0.75 value:2.5 > quantile:<quantile:0.9 value:2.5 > quantile:<quantile:0.95 value:2.5 > quantile:<quantile:0.99 value:2.5 > quantile:<quantile:0.999 value:2.5 > >",
		strings.TrimSpace(m.Metric[2].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"dog\" > summary:<sample_count:2 sample_sum:5.9 quantile:<quantile:0.5 value:2.6 > quantile:<quantile:0.75 value:3.3 > quantile:<quantile:0.9 value:3.3 > quantile:<quantile:0.95 value:3.3 > quantile:<quantile:0.99 value:3.3 > quantile:<quantile:0.999 value:3.3 > >",
		strings.TrimSpace(m.Metric[3].String()))
}

func TestSummaryWithLabels(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "BLAH"})

	s := metrics.SummaryWithLabels("animal breeds", []string{"animal", "breed"})
	s.Observe(1.0, "cat", "tabby")
	s.Observe(2.5, "cat", "siamese")
	s.Observe(2.6, "dog", "mutt")

	m := findMetric("blah_animal_breeds", metrics.gatherOK(t))
	assert.Equal(t, 3, len(m.Metric))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > label:<name:\"breed\" value:\"siamese\" > summary:<sample_count:1 sample_sum:2.5 quantile:<quantile:0.5 value:2.5 > quantile:<quantile:0.75 value:2.5 > quantile:<quantile:0.9 value:2.5 > quantile:<quantile:0.95 value:2.5 > quantile:<quantile:0.99 value:2.5 > quantile:<quantile:0.999 value:2.5 > >",
		strings.TrimSpace(m.Metric[0].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"cat\" > label:<name:\"breed\" value:\"tabby\" > summary:<sample_count:1 sample_sum:1 quantile:<quantile:0.5 value:1 > quantile:<quantile:0.75 value:1 > quantile:<quantile:0.9 value:1 > quantile:<quantile:0.95 value:1 > quantile:<quantile:0.99 value:1 > quantile:<quantile:0.999 value:1 > >",
		strings.TrimSpace(m.Metric[1].String()))
	assert.Equal(t, "label:<name:\"animal\" value:\"dog\" > label:<name:\"breed\" value:\"mutt\" > summary:<sample_count:1 sample_sum:2.6 quantile:<quantile:0.5 value:2.6 > quantile:<quantile:0.75 value:2.6 > quantile:<quantile:0.9 value:2.6 > quantile:<quantile:0.95 value:2.6 > quantile:<quantile:0.99 value:2.6 > quantile:<quantile:0.999 value:2.6 > >",
		strings.TrimSpace(m.Metric[2].String()))
}

func TestRegisterUnderlyingMetric(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "unused"})

	metricName := "name_" + fmt.Sprint(rand.Intn(100000))
	metric := prometheus.NewCounter(prometheus.CounterOpts{Name: metricName, Help: "help"})
	metrics.RegisterMetric(metric)
	metric.Add(71)

	m := findMetric(metricName, metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "counter:<value:71 >", strings.TrimSpace(m.Metric[0].String()))
}

func TestHistogramForResponseTime(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "A"})

	s := metrics.HistogramForResponseTime("MyHisto")
	s.Update(1.3)
	s.Update(2.5)
	s.Update(2.6)
	s.Update(2.9)
	s.Update(3.2)
	s.Update(3.3)
	s.Update(3.834344)

	m := findMetric("a_myhisto", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "histogram:<sample_count:7 sample_sum:19.634344000000002 bucket:<cumulative_count:0 upper_bound:0.005 > bucket:<cumulative_count:0 upper_bound:0.01 > bucket:<cumulative_count:0 upper_bound:0.025 > bucket:<cumulative_count:0 upper_bound:0.05 > bucket:<cumulative_count:0 upper_bound:0.1 > bucket:<cumulative_count:0 upper_bound:0.25 > bucket:<cumulative_count:0 upper_bound:0.5 > bucket:<cumulative_count:0 upper_bound:1 > bucket:<cumulative_count:2 upper_bound:2.5 > bucket:<cumulative_count:7 upper_bound:5 > bucket:<cumulative_count:7 upper_bound:10 > >",
		strings.TrimSpace(m.Metric[0].String()))
}

func TestHistogramCustomBuckets(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "A"})

	s := metrics.Histogram("MyHisto", []float64{2.0, 3.0, 3.5})
	s.Update(1.3)
	s.Update(2.5)
	s.Update(2.6)
	s.Update(2.9)
	s.Update(3.2)
	s.Update(3.3)
	s.Update(3.834344)

	m := findMetric("a_myhisto", metrics.gatherOK(t))
	assert.Equal(t, 1, len(m.Metric))
	assert.Equal(t, "histogram:<sample_count:7 sample_sum:19.634344000000002 bucket:<cumulative_count:1 upper_bound:2 > bucket:<cumulative_count:4 upper_bound:3 > bucket:<cumulative_count:6 upper_bound:3.5 > >",
		strings.TrimSpace(m.Metric[0].String()))
}

func TestErrors(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "z"})
	doTestErrors(t, &metrics)
	assert.Equal(t, []string{"z_errors"}, metrics.TestHelper().MetricNames())
}

func doTestErrors(t *testing.T, metrics *PrometheusMetricsImpl) {
	metrics.Error("bad")
	metrics.Error("generic")
	metrics.Error("generic")
	metrics.Error("worse")

	m := findMetric("z_errors", metrics.gatherOK(t))
	assert.Equal(t, 3, len(m.Metric))
	assert.Equal(t, "label:<name:\"error_type\" value:\"bad\" > counter:<value:1 >", strings.TrimSpace(m.Metric[0].String()))
	assert.Equal(t, "label:<name:\"error_type\" value:\"generic\" > counter:<value:2 >", strings.TrimSpace(m.Metric[1].String()))
	assert.Equal(t, "label:<name:\"error_type\" value:\"worse\" > counter:<value:1 >", strings.TrimSpace(m.Metric[2].String()))
}

func TestClear(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "z"})
	doTestCounter(t, &metrics)
	metrics.TestHelper().Clear()
	doTestErrors(t, &metrics)
	metrics.TestHelper().Clear()
	doTestCounter(t, &metrics)
}

func TestNames(t *testing.T) {
	metrics := NewMetrics(MetricOpts{Registry: prometheus.NewRegistry(), MetricNamePrefix: "blah"})
	metrics.Counter("c")
	metrics.Error("e")
	metrics.Gauge("g")
	metrics.HistogramForResponseTime("h")
	metrics.Summary("s")
	timedMethod(&metrics)
	assert.ElementsMatch(t, []string{"blah_c", "blah_g", "blah_h", "blah_s", "blah_timer", "blah_errors"}, metrics.TestHelper().MetricNames())
}

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

func timedMethod(metrics PrometheusMetrics) {
	defer metrics.Timer("Timer")()
	fmt.Println("Whatever it is we're timing")
}

func timedMethodWithLabel(metrics PrometheusMetrics) {
	defer metrics.TimerWithLabel("animal_timer", "animal", "cat")()
	fmt.Println("Whatever it is we're timing")
}

func (p *PrometheusMetricsImpl) gatherOK(t *testing.T) []*clientmodel.MetricFamily {
	gathered, err := p.TestHelper().Gather()
	assert.Nil(t, err)
	return gathered
}

type controlledTimer struct {
	eventTimer
	observer    prometheus.Observer
	expectation time.Duration
}

type controlledTimerFactory struct {
	timerFactory
	defaultExpectation time.Duration
}

func (f *controlledTimerFactory) NewTimer(o prometheus.Observer) eventTimer {
	return &controlledTimer{observer: o, expectation: f.defaultExpectation}
}

func (t *controlledTimer) Observe() time.Duration {
	fmt.Println("Making fixed observation...")
	t.observer.Observe(t.expectation.Seconds())
	return t.expectation
}

func findMetric(metricName string, gathered []*clientmodel.MetricFamily) *clientmodel.MetricFamily {
	for _, b := range gathered {
		if metricName == *b.Name {
			return b
		}
	}
	return nil
}

func assertRegistryHasMetricWithConfigTypeCounts(t *testing.T, metricLabelMap LabelsMap, metricName string, expectedLabelMap map[string]map[string]float64) {
	for labelName, labelExpectations := range expectedLabelMap {
		if assert.Containsf(t, metricLabelMap, labelName, "Metric %s did not contain expected label name %s.", metricName, labelName) {
			for labelValue, expectedCount := range labelExpectations {
				if assert.Containsf(t, metricLabelMap[labelName], labelValue, "Metric %s had label name %s but not label value %s", metricName, labelName, labelValue) {
					actualCount := metricLabelMap[labelName][labelValue].GetCounter().GetValue()
					assert.Equal(t, expectedCount, actualCount, `Expected metric %s{"%s": "%s"} to equal %.2f, but was %.2f`, metricName, labelName, labelValue, expectedCount, actualCount)
				}
			}
		}
	}
}

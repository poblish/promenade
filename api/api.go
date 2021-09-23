package api

import (
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metricFacade interface{}

type MetricDescriptions map[string]string

type MetricOpts struct {
	Registry                 prometheus.Registerer
	MetricNamePrefix         string
	PrefixSeparator          string
	Descriptions             MetricDescriptions
	CaseSensitiveMetricNames bool // true is faster, default is Insensitive
}

type PrometheusMetrics interface {
	Counter(name string, optionalDesc ...string) CounterFacade
	CounterWithLabel(name string, labelName string, optionalDesc ...string) LabelledCounterFacade
	CounterWithLabels(name string, labelNames []string, optionalDesc ...string) LabelledCounterFacade
	Error(name string) ErrorCounter
	Gauge(name string, optionalDesc ...string) GaugeFacade
	GaugeWithLabel(name string, labelName string, optionalDesc ...string) LabelledGaugeFacade
	GaugeWithLabels(name string, labelNames []string, optionalDesc ...string) LabelledGaugeFacade
	Histogram(name string, buckets []float64, optionalDesc ...string) HistogramFacade
	HistogramForResponseTime(name string, optionalDesc ...string) HistogramFacade
	Summary(name string, optionalDesc ...string) SummaryFacade
	SummaryWithLabel(name string, labelName string, optionalDesc ...string) LabelledSummaryFacade
	SummaryWithLabels(name string, labelNames []string, optionalDesc ...string) LabelledSummaryFacade
	Timer(Name string) func() time.Duration
	TimerWithLabel(Name string, labelName string, labelValue string) func() time.Duration
}

type PrometheusMetricsImpl struct {
	registry         prometheus.Registerer
	metricNamePrefix string
	descriptions     MetricDescriptions
	errorCounter     *prometheus.CounterVec
	errorCounterName string
	registrations    MetricRegistrations
	timerFactory     timerFactory

	caseSensitiveMetricNames bool // true is faster, default is Insensitive
	normalisedNames          normalisedNames
}

func NewMetrics(opts MetricOpts) PrometheusMetricsImpl {
	if opts.PrefixSeparator == "" {
		opts.PrefixSeparator = "_" // as per Prometheus lib standard
	}

	prefix := NormaliseAndLowercaseName(opts.MetricNamePrefix)
	if prefix != "" && !strings.HasSuffix(prefix, opts.PrefixSeparator) {
		prefix += opts.PrefixSeparator
	}

	if opts.Registry == nil {
		opts.Registry = prometheus.DefaultRegisterer
	}

	return PrometheusMetricsImpl{registry: opts.Registry,
		metricNamePrefix:         prefix,
		descriptions:             opts.Descriptions,
		registrations:            newMetricRegistrations(),
		timerFactory:             &defaultTimerFactory{},
		caseSensitiveMetricNames: opts.CaseSensitiveMetricNames,
		normalisedNames:          normalisedNames{internal: make(map[string]string)},
	}
}

const (
	TypeCounter       = iota << 2
	TypeCounterLabels = iota << 2
	TypeGauge         = iota << 2
	TypeGaugeLabels   = iota << 2
	TypeSummary       = iota << 2
	TypeSummaryLabels = iota << 2
	TypeHistogram     = iota << 2
)

type metricEntry struct {
	metric     metricFacade
	metricType int
}

type MetricRegistrations struct {
	sync.RWMutex
	internal map[string]metricEntry
}

func newMetricRegistrations() MetricRegistrations {
	return MetricRegistrations{internal: make(map[string]metricEntry)}
}

type normalisedNames struct {
	sync.RWMutex
	internal map[string]string
}

func newNormalisedNames() normalisedNames {
	return normalisedNames{internal: make(map[string]string)}
}

func (p *PrometheusMetricsImpl) RegisterMetric(metric prometheus.Collector) {
	p.registry.Register(metric)
}

type MetricBuilder func(p *PrometheusMetricsImpl, name string, desc string) interface{}

func (p *PrometheusMetricsImpl) getOrAdd(name string, metricType int, builder MetricBuilder, desc []string) metricFacade {
	var metricKey string

	if p.caseSensitiveMetricNames {
		metricKey = normalizer.Replace(name)
	} else {
		if entry, ok := p.getNormalisedName(name); ok {
			metricKey = entry
		} else {
			metricKey = NormaliseAndLowercaseName(name)
			p.storeNormalisedName(name, metricKey)
		}
	}

	if entry, ok := p.getRegistration(metricKey); ok {
		if entry.metricType != metricType {
			panic(p.getFullMetricName(metricKey) + " is already used for a different type of metric")
		}
		return entry.metric
	}

	var newMetric = builder(p, p.getFullMetricName(metricKey), p.bestDescription(metricKey, desc))
	p.storeRegistration(metricKey, metricEntry{metric: newMetric, metricType: metricType})
	return newMetric
}

func (p *PrometheusMetricsImpl) getNormalisedName(name string) (string, bool) {
	p.normalisedNames.RLock()
	defer p.normalisedNames.RUnlock()
	val, ok := p.normalisedNames.internal[name]
	return val, ok
}
func (p *PrometheusMetricsImpl) storeNormalisedName(name string, value string) {
	p.normalisedNames.Lock()
	defer p.normalisedNames.Unlock()
	p.normalisedNames.internal[name] = value
}

func (p *PrometheusMetricsImpl) getRegistration(key string) (metricEntry, bool) {
	p.registrations.RLock()
	defer p.registrations.RUnlock()
	val, ok := p.registrations.internal[key]
	return val, ok
}

func (p *PrometheusMetricsImpl) storeRegistration(key string, value metricEntry) {
	p.registrations.Lock()
	defer p.registrations.Unlock()
	p.registrations.internal[key] = value
}

func (p *PrometheusMetricsImpl) bestDescription(name string, desc []string) string {
	var description = ""
	if len(desc) > 0 {
		description = desc[0]
	}

	if description == "" {
		if mapping, found := p.descriptions[name]; found {
			description = mapping
		} else {
			description = p.getFullMetricName(name)
		}
	}
	return description
}

func (p *PrometheusMetricsImpl) getFullMetricName(name string) string {
	return p.metricNamePrefix + name
}

var normalizer = strings.NewReplacer(".", "_", "-", "_", "#", "_", " ", "_")

func NormaliseAndLowercaseName(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToLower(normalizer.Replace(name))
}

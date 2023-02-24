package test

import "github.com/poblish/promenade/api"

type emptyCounter struct {
}

func (e emptyCounter) Inc() {
	// NOOP
}

func (e emptyCounter) IncBy(inc float64) {
	// NOOP
}

type emptyLabelledCounter struct {
}

func (e emptyLabelledCounter) IncLabelBy(labelValues ...string) api.IncByValue {
	return api.IncByValue{}
}

func (e emptyLabelledCounter) IncLabel(labelValues ...string) {
	// NOOP
}

func (e emptyLabelledCounter) Value(inc float64) {
	// NOOP
}

type emptyGauge struct {
}

func (e emptyGauge) SetValue(value float64) {
	// NOOP
}

func (e emptyGauge) Inc() {
	// NOOP
}

func (e emptyGauge) IncBy(inc float64) {
	// NOOP
}

func (e emptyGauge) Dec() {
	// NOOP
}

func (e emptyGauge) DecBy(dec float64) {
	// NOOP
}

type emptyHistogram struct {
}

func (e emptyHistogram) Update(value float64) {
	// NOOP
}

type emptyLabelledGauge struct {
}

func (e emptyLabelledGauge) IncLabels(labelValues ...string) {
	// NOOP
}

func (e emptyLabelledGauge) IncLabelsBy(labelValues ...string) api.IncGaugeByValue {
	return api.IncGaugeByValue{}
}

func (e emptyLabelledGauge) DecLabels(labelValues ...string) {
	// NOOP
}

func (e emptyLabelledGauge) DecLabelsBy(labelValues ...string) api.DecGaugeByValue {
	return api.DecGaugeByValue{}
}

func (e emptyLabelledGauge) SetLabels(labelValues ...string) api.SetGaugeByValue {
	return api.SetGaugeByValue{}
}

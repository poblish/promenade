package test

import "github.com/poblish/promenade/api"

type emptyCounter struct {
}

func (e emptyCounter) Inc() {
	// NOOP
}

func (e emptyCounter) IncBy(_ float64) {
	// NOOP
}

type emptyLabelledCounter struct {
}

type emptyIncByValue struct {
}

func (e emptyIncByValue) Value(_ float64) {
	// NOOP
}

func (e emptyLabelledCounter) IncLabelBy(_ ...string) api.IncByValue {
	return emptyIncByValue{}
}

func (e emptyLabelledCounter) IncLabel(_ ...string) {
	// NOOP
}

func (e emptyLabelledCounter) Value(_ float64) {
	// NOOP
}

type emptyGauge struct {
}

func (e emptyGauge) SetValue(_ float64) {
	// NOOP
}

func (e emptyGauge) Inc() {
	// NOOP
}

func (e emptyGauge) IncBy(_ float64) {
	// NOOP
}

func (e emptyGauge) Dec() {
	// NOOP
}

func (e emptyGauge) DecBy(_ float64) {
	// NOOP
}

type emptyHistogram struct {
}

func (e emptyHistogram) Update(_ float64) {
	// NOOP
}

type emptyLabelledGauge struct {
}

func (e emptyLabelledGauge) IncLabels(_ ...string) {
	// NOOP
}

type emptyIncGaugeByValue struct {
}

func (e emptyIncGaugeByValue) Value(inc float64) {
	// NOOP
}

func (e emptyLabelledGauge) IncLabelsBy(_ ...string) api.IncGaugeByValue {
	return emptyIncGaugeByValue{}
}

func (e emptyLabelledGauge) DecLabels(_ ...string) {
	// NOOP
}

type emptyDecGaugeByValue struct {
}

func (e emptyDecGaugeByValue) Value(inc float64) {
	// NOOP
}

func (e emptyLabelledGauge) DecLabelsBy(_ ...string) api.DecGaugeByValue {
	return emptyDecGaugeByValue{}
}

type emptySetGaugeByValue struct {
}

func (e emptySetGaugeByValue) Value(inc float64) {
	// NOOP
}

func (e emptyLabelledGauge) SetLabels(_ ...string) api.SetGaugeByValue {
	return emptySetGaugeByValue{}
}

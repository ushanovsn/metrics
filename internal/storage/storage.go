package storage

import (
	"fmt"
	"sort"
)

// storage of metrics
type MemStorage struct {
	metrics Metrics
}

// metrics values
type Metrics struct {
	// values of gauge
	Gauge map[string]float64
	// values of counter
	Counter map[string]int64
}

// storage interface
type Repositories interface {
	SetGauge(string, float64)
	SetCounter(string, int64)
	GetGauge(string) (float64, bool)
	GetCounter(string) (int64, bool)
}

// init metric storage (at this time)
var Metr MemStorage = MemStorage{
	metrics: Metrics{
		Gauge: make(map[string]float64),
		Counter: make(map[string]int64),
	},
}


// add or update of gauge metric
func (ms *MemStorage) SetGauge(name string, val float64) {
	ms.metrics.Gauge[name] = val
}


// add or update of counter metric
func (ms *MemStorage) SetCounter(name string, val int64) {
	ms.metrics.Counter[name] += val
}

// get of gauge metric
func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	v, ok := ms.metrics.Gauge[name]
	return v, ok
}

// get of gauge metrics list
func (ms *MemStorage) GetGaugeList() []string {
	var list []string
	for k, v := range ms.metrics.Gauge {
		list = append(list, fmt.Sprintf("Name: %-10s,\tValue: %v", k, v))
	}
	sort.Strings(sort.StringSlice(list))
	return list
}


// get of counter metric
func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	v, ok := ms.metrics.Counter[name]
	return v, ok
}

// get of counter metrics list
func (ms *MemStorage) GetCounterList() []string {
	var list []string
	for k, v := range ms.metrics.Counter {
		list = append(list, fmt.Sprintf("Name: %-10s,\tValue: %v", k, v))
	}
	sort.Strings(sort.StringSlice(list))
	return list
}

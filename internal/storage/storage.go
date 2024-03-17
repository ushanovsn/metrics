package storage

import (
	"fmt"
	"sort"
)

// storage of metrics
type MemStorage struct {
	metrics metrics
}

// metrics values
type metrics struct {
	// values of gauge
	gauge map[string]float64
	// values of counter
	counter map[string]int64
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
	metrics: metrics{
		gauge: make(map[string]float64),
		counter: make(map[string]int64),
	},
}


// add or update of gauge metric
func (ms *MemStorage) SetGauge(name string, val float64) {
	ms.metrics.gauge[name] = val
}


// add or update of counter metric
func (ms *MemStorage) SetCounter(name string, val int64) {
	ms.metrics.counter[name] += val
}

// get of gauge metric
func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	v, ok := ms.metrics.gauge[name]
	return v, ok
}

// get of gauge metrics list
func (ms *MemStorage) GetGaugeList() []string {
	list := []string{}
	for k, v := range ms.metrics.gauge {
		list = append(list, fmt.Sprintf("Name: %s,\tValue: %v", k, v))
	}
	sort.Strings(sort.StringSlice(list))
	return list
}


// get of counter metric
func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	v, ok := ms.metrics.counter[name]
	return v, ok
}

// get of counter metrics list
func (ms *MemStorage) GetCounterList() []string {
	list := []string{}
	for k, v := range ms.metrics.counter {
		list = append(list, fmt.Sprintf("Name: %s,\tValue: %v", k, v))
	}
	sort.Strings(sort.StringSlice(list))
	return list
}

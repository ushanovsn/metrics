package storage


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

package metricsproc

import (
	"runtime"
	"reflect"
	"math/rand"
)

// type for metrics param
type Metrics struct {
	TypeM	string
	ValueF 	float64
	ValueI 	int64
}

// init storage for actual metrics
var MetrStor map[string]map[string]Metrics = map[string]map[string]Metrics{ 
	"gauge": make(map[string]Metrics),
	"counter": make(map[string]Metrics),
}


// main func for collecting and updating metrics
func MetrCollect() {
	runtimeCollect()
	counterCollect()
	randomCollect()
}


// updating runtime metrics
func runtimeCollect() {
	var rt runtime.MemStats
	//read metrics
	runtime.ReadMemStats(&rt)

    v := reflect.ValueOf(rt)
    for i := 0; i < v.NumField(); i++ {
		if en, ok := mNames[v.Type().Field(i).Name]; ok && en {
			// ckeck type of metric value
			switch v.Type().Field(i).Type.String() {
			case "uint64", "uint32":
				// set type of this metric and value
				MetrStor["gauge"][v.Type().Field(i).Name] = Metrics{TypeM: "float64", ValueF: float64(v.Field(i).Uint())}
			case "int64", "int32":
				// set type of this metric and value
				MetrStor["gauge"][v.Type().Field(i).Name] = Metrics{TypeM: "float64", ValueF: float64(v.Field(i).Int())}
			case "float64", "float32":
				// set type of this metric and value
				MetrStor["gauge"][v.Type().Field(i).Name] = Metrics{TypeM: "float64", ValueF: v.Field(i).Float()}
			}
		}
    }
}



// updating counter metrics
func counterCollect() {
	if m, ok := MetrStor["counter"]["PollCount"]; ok {
		// add count
		MetrStor["counter"]["PollCount"] = Metrics{TypeM: m.TypeM, ValueI: (m.ValueI + 1)}
	} else {
		// init metric
		MetrStor["counter"]["PollCount"] = Metrics{TypeM: "int64", ValueI: 1}
	}
}


// updating random gauge metrics
func randomCollect() {
	val := rand.Float64()
	if m, ok := MetrStor["gauge"]["RandomValue"]; ok {
		// update value
		MetrStor["gauge"]["RandomValue"] = Metrics{TypeM: m.TypeM, ValueF: val}
	} else {
		// init metric
		MetrStor["gauge"]["RandomValue"] = Metrics{TypeM: "float64", ValueF: val}
	}
}
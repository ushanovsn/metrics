package metricscollector

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
)

// type for metrics param
type metrics struct {
	typeM	string
	valueF 	float64
	valueI 	int64
}

// init storage for actual metrics
var MetrStor map[string]map[string]metrics = map[string]map[string]metrics{ 
	"gauge": make(map[string]metrics),
	"counter": make(map[string]metrics),
}


func GetMetricsList() []string {
	list := make([]string, 0, len(mGaugeNames) + len(mCounterNames) + len(mRandomNames))
	var elem strings.Builder

	for t, v := range MetrStor {
		// check metrics
		for n, m := range v {
			elem.WriteString("/" + t + "/" + n)
			switch m.typeM {
			case "float64":
				elem.WriteString("/" + fmt.Sprint(m.valueF))
			case "int64":
				elem.WriteString("/" + fmt.Sprint(m.valueI))
			default:
				elem.WriteString("/")
			}

			list = append(list, elem.String())
		}
	}
	return list
}

// main func for collecting and updating metrics
func MetrCollect() {
	runtimeCollect()
	counterCollect()
	randomCollect()
}


// updating gauge runtime metrics
func runtimeCollect() {
	var rt runtime.MemStats
	//read metrics
	runtime.ReadMemStats(&rt)

	v := reflect.ValueOf(rt)
	for i := 0; i < v.NumField(); i++ {
		if en, ok := mGaugeNames[v.Type().Field(i).Name]; ok && en {
			// ckeck type of metric value
			switch v.Type().Field(i).Type.String() {
			case "uint64", "uint32":
				// set type of this metric and value
				MetrStor["gauge"][v.Type().Field(i).Name] = metrics{typeM: "float64", valueF: float64(v.Field(i).Uint())}
			case "int64", "int32":
				// set type of this metric and value
				MetrStor["gauge"][v.Type().Field(i).Name] = metrics{typeM: "float64", valueF: float64(v.Field(i).Int())}
			case "float64", "float32":
				// set type of this metric and value
				MetrStor["gauge"][v.Type().Field(i).Name] = metrics{typeM: "float64", valueF: v.Field(i).Float()}
			}
		}
	}
}



// updating counter metrics
func counterCollect() {
	for n, f := range mCounterNames {
		if f {
			if m, ok := MetrStor["counter"][n]; ok {
				// add count
				MetrStor["counter"][n] = metrics{typeM: m.typeM, valueI: (m.valueI + 1)}
			} else {
				// init metric
				MetrStor["counter"][n] = metrics{typeM: "int64", valueI: 1}
			}
		}
	}
}


// updating gauge random metrics
func randomCollect() {
	for n, f := range mRandomNames {
		if f {
			val := rand.Float64()
			if m, ok := MetrStor["gauge"][n]; ok {
				// update value
				MetrStor["gauge"][n] = metrics{typeM: m.typeM, valueF: val}
			} else {
				// init metric
				MetrStor["gauge"][n] = metrics{typeM: "float64", valueF: val}
			}
		}
	}
}
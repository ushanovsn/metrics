package agent

import (
	"log"
	"fmt"
	"github.com/ushanovsn/metrics/internal/options"
	"net/http"
	"time"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
)




func AgentRun() error {
	AgentOpt := options.AgentOptions{
		Net: options.NetAddress{
			Host: "localhost",
			Port: 8080,
		},
		ReportInterval: 10,
		PollInterval: 2,
	}

	var err error
	var timer int
	var sleepTime int

	var minTimer int
	var maxTimer int

	InitFlag(&AgentOpt)
	InitEnv(&AgentOpt)

	actualValG := initGaugeValues()

	actualValC := initCounterValues()

	repInt := AgentOpt.GetRepInt()
	pollInt := AgentOpt.GetPollInt()

	log.Printf("Agent run with poll interval %v and report interval %v\n", pollInt, repInt)

	if repInt < pollInt {
		minTimer = repInt
		maxTimer = pollInt
	} else {
		minTimer = pollInt
		maxTimer = repInt
	}

	for {

		switch {
		case minTimer == maxTimer:
			sleepTime = minTimer
			timer = sleepTime
		case (timer < minTimer) || (timer < maxTimer && (maxTimer-timer) >= minTimer):
			// first start, next minTimer iteration when are multiples
			sleepTime = minTimer
			timer += sleepTime
		case timer < maxTimer && (maxTimer-timer) < minTimer:
			// last iteration when minTimer and maxTimer is not multiples
			sleepTime = (maxTimer - timer)
			timer += sleepTime
		case (timer >= maxTimer):
			// after maxTimer iteration
			sleepTime = minTimer
			timer = sleepTime
		}

		// waiting...
		time.Sleep(time.Duration(sleepTime) * time.Second)

		if minTimer == maxTimer ||(timer == maxTimer && maxTimer % minTimer == 0) {
			getGauge(actualValG)
			getCounter(actualValC)
			if err != nil {
				return err
			}

			err = agentSendMetrics(actualValG, actualValC, AgentOpt.Net.String())
			if err != nil {
				return err
			}

		} else {
			if timer >= minTimer && timer < maxTimer {
				// tick of minTimer
				if minTimer == pollInt {
					getGauge(actualValG)
					getCounter(actualValC)
				} else {
					err = agentSendMetrics(actualValG, actualValC, AgentOpt.Net.String())
				}
			} else {
				// tick of maxTimer
				if maxTimer == pollInt {
					getGauge(actualValG)
					getCounter(actualValC)
				} else {
					err = agentSendMetrics(actualValG, actualValC, AgentOpt.Net.String())
				}
			}

			if err != nil {
				return err
			}
		}
	}
}



func agentSendMetrics(gm *map[string]float64, cm *map[string]int64, addr string) error {
	client := &http.Client{}

	//log.Printf("send: %v\n", metricscollector.GetMetricsList())

	for _, val := range metricsToList(gm, cm) {

		log.Printf("* send path: %s\n", val)

		// POST request with metric
		r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update%s", addr, val), nil)
		if err != nil {
			return err
		}
		

		// add header (optional)
		r.Header.Add("Content-Type", "text/plain")

		// execute POST request
		resp, err := client.Do(r)

		if err != nil {
			log.Printf("error while requesting: \"%s\"\n", err)
		} else {
			resp.Body.Close()
		}
		
	}
	return nil
}



func getGauge(gm *map[string]float64) {
	var rt runtime.MemStats
	//read metrics
	runtime.ReadMemStats(&rt)

	v := reflect.ValueOf(rt)
	for i := 0; i < v.NumField(); i++ {
		if _, ok := (*gm)[v.Type().Field(i).Name]; ok {
			switch v.Type().Field(i).Type.String() {
			case "uint64", "uint32":
				// set type of this metric and value
				(*gm)[v.Type().Field(i).Name] = float64(v.Field(i).Uint())
			case "int64", "int32":
				// set type of this metric and value
				(*gm)[v.Type().Field(i).Name] = float64(v.Field(i).Int())
			case "float64", "float32":
				// set type of this metric and value
				(*gm)[v.Type().Field(i).Name] = v.Field(i).Float()
			}
		}
	}

	(*gm)["RandomValue"] =  rand.Float64()
}

func getCounter(cm *map[string]int64) {
	for i := range *cm {
		(*cm)[i]++
	}
}


func initGaugeValues() *map[string]float64 {
	metricsNames := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC",  "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",}
	gm := make(map[string]float64)
	for _, n := range metricsNames {
		gm[n] = 0.0
	}
	return &gm
}

func initCounterValues() *map[string]int64 {
	metricsNames := []string{"PollCount",}
	cm := make(map[string]int64)
	for _, n := range metricsNames {
		cm[n] = 0.0
	}
	return &cm
}


func metricsToList(gm *map[string]float64, cm *map[string]int64) []string {
	list := make([]string, 0, len(*gm) + len(*cm))
	var elem strings.Builder

	for t, v := range *gm {
		elem.WriteString("/" + t + "/" + fmt.Sprint(v))
		list = append(list, elem.String())
	}
	

	for t, v := range *cm {
		elem.WriteString("/" + t + "/" + fmt.Sprint(v))
		list = append(list, elem.String())
	}
	
	return list
}


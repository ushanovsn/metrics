package agent

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/logger"
	"github.com/ushanovsn/metrics/internal/options"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type agentData struct {
	opt options.AgentOptions
}

func AgentInit() *agentData {
	agentOpt := options.InitAg()

	InitFlag(agentOpt)
	envErr := InitEnv(agentOpt)

	logger.InitLogger(&agentOpt.Logger)
	log := agentOpt.Logger.GetLogger()

	if envErr != nil {
		log.Errorf("Error when parsing environment var: %s\n", envErr)
	}

	return &agentData{
		opt: *agentOpt,
	}
}

func AgentRun(ag *agentData) error {

	var err error
	var timer int
	var sleepTime int

	var minTimer int
	var maxTimer int

	log := ag.opt.Logger.GetLogger()

	actualValG := initGaugeValues()

	actualValC := initCounterValues()

	repInt := ag.opt.GetRepInt()
	pollInt := ag.opt.GetPollInt()

	log.Debugf("Agent run with poll interval %v and report interval %v\n", pollInt, repInt)

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

		if minTimer == maxTimer || (timer == maxTimer && maxTimer%minTimer == 0) {
			getGauge(actualValG)
			getCounter(actualValC)
			if err != nil {
				return err
			}

			err = agentSendMetrics(actualValG, actualValC, ag)
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
					err = agentSendMetrics(actualValG, actualValC, ag)
				}
			} else {
				// tick of maxTimer
				if maxTimer == pollInt {
					getGauge(actualValG)
					getCounter(actualValC)
				} else {
					err = agentSendMetrics(actualValG, actualValC, ag)
				}
			}

			if err != nil {
				return err
			}
		}
	}
}

func agentSendMetrics(gm *map[string]float64, cm *map[string]int64, ag *agentData) error {
	client := &http.Client{}
	log := ag.opt.Logger.GetLogger()

	for _, val := range metricsToList(gm, cm) {

		log.Debugf("* POST: %s\n", fmt.Sprintf("http://%s/update%s", ag.opt.Net.String(), val))

		// POST request with metric
		r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update%s", ag.opt.Net.String(), val), nil)
		if err != nil {
			return err
		}

		// add header (optional)
		r.Header.Add("Content-Type", "text/plain")

		// execute POST request
		resp, err := client.Do(r)

		if err != nil {
			log.Errorf("Error while POST executing: \"%s\"\n", err)
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

	(*gm)["RandomValue"] = rand.Float64()
}

func getCounter(cm *map[string]int64) {
	for i := range *cm {
		(*cm)[i]++
	}
}

func initGaugeValues() *map[string]float64 {
	metricsNames := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}
	gm := make(map[string]float64)
	for _, n := range metricsNames {
		gm[n] = 0.0
	}
	return &gm
}

func initCounterValues() *map[string]int64 {
	metricsNames := []string{"PollCount"}
	cm := make(map[string]int64)
	for _, n := range metricsNames {
		cm[n] = 0.0
	}
	return &cm
}

func metricsToList(gm *map[string]float64, cm *map[string]int64) []string {
	list := make([]string, 0, len(*gm)+len(*cm))
	var elem strings.Builder

	for t, v := range *gm {
		elem.WriteString("/gauge/" + t + "/" + fmt.Sprint(v))
		list = append(list, elem.String())
		elem.Reset()
	}

	for t, v := range *cm {
		elem.WriteString("/counter/" + t + "/" + fmt.Sprint(v))
		list = append(list, elem.String())
		elem.Reset()
	}

	return list
}

func AgentStop(ag *agentData) {
	err := ag.opt.Logger.Stop()
	if err != nil {
		fmt.Printf("Error while stopping Agent: %s\n", err)
	}
}

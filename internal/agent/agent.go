package agent

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/metricsproc"
	"net/http"
	"time"
)

var client *http.Client

func AgentRun() error {

	client = &http.Client{}

	var err error
	var timer int
	var sleepTime int

	var minTimer int
	var maxTimer int

	if FlagParam.ReportInterval < FlagParam.PollInterval {
		minTimer = FlagParam.ReportInterval
		maxTimer = FlagParam.PollInterval
	} else {
		minTimer = FlagParam.PollInterval
		maxTimer = FlagParam.ReportInterval
	}

	fmt.Printf("FlagParam: %v\n", FlagParam)
	fmt.Printf("minTimer: %v, maxTimer: %v\n", minTimer, maxTimer)

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

		//fmt.Printf("minTimer: %v, maxTimer: %v\n", minTimer, maxTimer)
		fmt.Printf("sleepTime: %v, timer: %v\n", sleepTime, timer)

		// waiting...
		time.Sleep(time.Duration(sleepTime) * time.Second)

		if minTimer == maxTimer ||(timer == maxTimer && maxTimer % minTimer == 0) {
			err = AgentCheckMetrics()
			if err != nil {
				return err
			}

			err = AgentSendMetrics()
			if err != nil {
				return err
			}

		} else {
			if timer >= minTimer && timer < maxTimer {
				// tick of minTimer
				if minTimer == FlagParam.PollInterval {
					err = AgentCheckMetrics()
				} else {
					err = AgentSendMetrics()
				}
			} else {
				// tick of maxTimer
				if maxTimer == FlagParam.PollInterval {
					err = AgentCheckMetrics()
				} else {
					err = AgentSendMetrics()
				}
			}

			if err != nil {
				return err
			}
		}
	}
}

func AgentCheckMetrics() error {

	// updating metrics
	metricsproc.MetrCollect()

	fmt.Printf("checked!\n")
	return nil
}

func AgentSendMetrics() error {

	for t, v := range metricsproc.MetrStor {
		// check metrics
		for n, m := range v {
			postPath := "/" + t + "/" + n
			switch m.TypeM {
			case "float64":
				postPath += "/" + fmt.Sprint(m.ValueF)
			case "int64":
				postPath += "/" + fmt.Sprint(m.ValueI)
			default:
				postPath += "/"
			}

			//fmt.Printf("* send path: %s\n", postPath)

			// POST request with metric
			r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update%s", FlagParam.Net.String(), postPath), nil)
			if err != nil {
				panic(err)
			}

			// add header (optional)
			r.Header.Add("Content-Type", "text/plain")

			// execute POST request
			resp, err := client.Do(r)
			if err != nil {
				fmt.Printf("error while requesting: %s \n", err)
			}

			resp.Body.Close()
		}

	}
	fmt.Printf("sended!\n")
	return nil
}

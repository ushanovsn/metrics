package agent

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/metricsproc"
	"github.com/ushanovsn/metrics/internal/options"
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

	if options.AgentOpt.ReportInterval < options.AgentOpt.PollInterval {
		minTimer = options.AgentOpt.ReportInterval
		maxTimer = options.AgentOpt.PollInterval
	} else {
		minTimer = options.AgentOpt.PollInterval
		maxTimer = options.AgentOpt.ReportInterval
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
				if minTimer == options.AgentOpt.PollInterval {
					err = AgentCheckMetrics()
				} else {
					err = AgentSendMetrics()
				}
			} else {
				// tick of maxTimer
				if maxTimer == options.AgentOpt.PollInterval {
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

			fmt.Printf("* send path: %s\n", postPath)

			// POST request with metric
			r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update%s", options.AgentOpt.Net.String(), postPath), nil)
			if err != nil {
				return err
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
	return nil
}

package agent

import (
	"log"
	"fmt"
	"github.com/ushanovsn/metrics/internal/metricscollector"
	"github.com/ushanovsn/metrics/internal/options"
	"net/http"
	"time"
)


func AgentRun() error {

	var err error
	var timer int
	var sleepTime int

	var minTimer int
	var maxTimer int

	repInt := options.AgentOpt.GetRepInt()
	pollInt := options.AgentOpt.GetPollInt()

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
			metricscollector.MetrCollect()
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
				if minTimer == pollInt {
					metricscollector.MetrCollect()
				} else {
					err = AgentSendMetrics()
				}
			} else {
				// tick of maxTimer
				if maxTimer == pollInt {
					metricscollector.MetrCollect()
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



func AgentSendMetrics() error {
	client := &http.Client{}

	log.Printf("send: %v\n", metricscollector.GetMetricsList())

	for _, val := range metricscollector.GetMetricsList() {

		log.Printf("* send path: %s\n", val)

		// POST request with metric
		r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update%s", options.AgentOpt.Net.String(), val), nil)
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

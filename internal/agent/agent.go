package agent

import (
	"fmt"
	"time"
	"net/http"
	"github.com/ushanovsn/metrics/internal/metricsproc"
)

const (
	pollInterval = 2
	reportInterval = 10
)


func StartAgent() {

	client := &http.Client{}

	var timer int

	for {
		for timer < reportInterval {

			// updating metrics
			metricsproc.MetrCollect()

			time.Sleep(pollInterval * time.Second)
			timer += pollInterval

			if timer >= reportInterval {
				timer = 0
				break
			}
		}

		// check types of metrics
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
				r, err := http.NewRequest("POST", "http://localhost:8080/update" + postPath, nil)
				if err != nil {
					panic(err)
				}
		
				// add header (optional)
				r.Header.Add("Content-Type", "text/plain")

				// execute POST request
				_, err = client.Do(r)
				if err != nil {
					fmt.Printf("error while requesting: %s \n", err)
				}
			}

		}

	}
}
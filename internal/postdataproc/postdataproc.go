package postdataproc

import (
	"github.com/ushanovsn/metrics/internal/storage"
	//"fmt"
	"strings"
	"strconv"
)

// error enums
type ProcError int

const (
	ProcNoErrors ProcError = iota
	ProcWrongType
	ProcWrongValue
	ProcWrongName
)

// process data, what received by HTTP POST 
func UsePOSTData(post []string) ProcError {
	procErr := ProcNoErrors
	
	if len(post) == 3 {
		mType := post[0]
		mName := post[1]
		mVal := post[2]

		// check metric name
		if strings.TrimSpace(mName) == "" {
			procErr = ProcWrongName
		} else {

			switch strings.ToLower(mType) {
			case "gauge":
				if v, err := strconv.ParseFloat(mVal, 64); err == nil {
					storage.Metr.SetGauge(mName, v)
					//fmt.Printf("Received gauge: %s, value: %v\n", mName, v)
				} else {
					procErr = ProcWrongValue
				}
			case "counter":
				if v, err := strconv.ParseInt(mVal, 10, 64); err == nil {
					storage.Metr.SetCounter(mName, v)
					//fmt.Printf("Received counter: %s, value: %v\n", mName, v)
				} else {
					procErr = ProcWrongValue
				}
			default:
				procErr = ProcWrongType
			}
		
		}

	} else if len(post) == 2 {
		// not the full number of parameters, we believe that we did not get the name of the metric
		procErr = ProcWrongName
	} else {
		// the values are not correct, the quantity does not correspond to the desired one
		procErr = ProcWrongValue
	}

	return procErr
}

package handlers

import (
	"net/http"
	"strings"
	"github.com/ushanovsn/metrics/internal/postdataproc"
	"fmt"
)


// start page if it needed
func startPage(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	//res.Write([]byte("Metrics holder server"))
}


// processing all received data by "update" address
func updatePage(res http.ResponseWriter, req *http.Request) {
	// post message if needed
	var msg []byte
	// default header = badRequest 
	header := http.StatusBadRequest

	// check method
	if req.Method == http.MethodPost {

		// check content type
		if true /* req.Header.Values("Content-Type") != nil && req.Header.Values("Content-Type")[0] == "text/plain" */ {
			// processing received data
			err := postdataproc.UsePOSTData(strings.Split(req.URL.Path, "/"))
			// check processing errors
			if err == postdataproc.ProcNoErrors {
				header = http.StatusOK
			} else {
				switch err {
				case postdataproc.ProcWrongType:
					msg = []byte("Wrong metric type")
				case postdataproc.ProcWrongValue:
					msg = []byte("Wrong metric value")
				case postdataproc.ProcWrongName:
					header = http.StatusNotFound
					msg = []byte("Can't find metric name")
				default:
					msg = []byte("Error occurred")
				}
			}
		} else {
			msg = []byte("Only \"text/plain\" accepted")
		}

	} else {
		msg = []byte("Only POST accepted")
	}


	res.WriteHeader(header)
	fmt.Printf("http msg: %s\n", msg)
}
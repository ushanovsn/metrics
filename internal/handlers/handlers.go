package handlers

import (
	"net/http"
	"strings"
	"github.com/ushanovsn/metrics/internal/rcvddataproc"
	"fmt"
)


func ServerMux() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", startPage)

	mux.Handle("/update/", http.StripPrefix("/update/", http.HandlerFunc(updatePage)))



	return mux
}

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
		var rightContentT = true
		for i, v := range req.Header {
			if i == "Content-Type" && v[0] != "text/plain" {
				rightContentT = false
			}
		}


		if rightContentT {
			// processing received data
			err := rcvddataproc.UsePOSTData(strings.Split(req.URL.Path, "/"))
			// check processing errors
			if err == rcvddataproc.ProcNoErrors {
				header = http.StatusOK
			} else {
				switch err {
				case rcvddataproc.ProcWrongType:
					msg = []byte("Wrong metric type")
				case rcvddataproc.ProcWrongValue:
					msg = []byte("Wrong metric value")
				case rcvddataproc.ProcWrongName:
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
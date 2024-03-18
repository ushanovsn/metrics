package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ushanovsn/metrics/internal/rcvddataproc"
	"github.com/ushanovsn/metrics/internal/storage"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// start page
func StartPage(res http.ResponseWriter, req *http.Request) {
	header := http.StatusOK

	res.Header().Add("Content-Type", "text/html")

	data := struct {
		Title    string
		MetricsG []string
		MetricsC []string
	}{
		Title:    "Metrics list",
		MetricsG: storage.Metr.GetGaugeList(),
		MetricsC: storage.Metr.GetCounterList(),
	}

	tmpl, err := template.ParseFiles("./static/htmltemplates/main_page_template.html")
	if err != nil {
		fmt.Printf("error while loading template: %v\n", err)
		header = http.StatusInternalServerError
		res.WriteHeader(header)
		return
	}

	res.WriteHeader(header)
	err = tmpl.Execute(res, data)

	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}

// metric page
func GetPageM(res http.ResponseWriter, req *http.Request) {
	var msg []byte
	header := http.StatusNotFound
	res.Header().Add("Content-Type", "text/plain")

	if chi.URLParam(req, "mName") == "" {
		res.WriteHeader(header)
		return
	}

	switch strings.ToLower(chi.URLParam(req, "mType")) {
	case "gauge":
		if v, ok := storage.Metr.GetGauge(chi.URLParam(req, "mName")); ok {
			header = http.StatusOK
			msg = []byte(fmt.Sprint(v))
		}
	case "counter":
		if v, ok := storage.Metr.GetCounter(chi.URLParam(req, "mName")); ok {
			header = http.StatusOK
			msg = []byte(fmt.Sprint(v))
		}
	}

	res.WriteHeader(header)
	if len(msg) > 0 {
		if _, err := res.Write(msg); err != nil {
			log.Printf("Error while write msg: %s\n", err)
		}

	}
}

// processing all received data by "update" address
func UpdatePageM(res http.ResponseWriter, req *http.Request) {
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
			rcvdData := []string{
				chi.URLParam(req, "mType"),
				chi.URLParam(req, "mName"),
				chi.URLParam(req, "mValue"),
			}

			// processing received data
			err := rcvddataproc.UsePOSTData(rcvdData)
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
	if len(msg) > 0 {
		fmt.Printf("http msg: %s\n", msg)
	}
}

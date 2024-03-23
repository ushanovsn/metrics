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
func StartPage(repo *storage.Repositories) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title    string
			MetricsG []string
			MetricsC []string
		}{
			Title:    "Metrics list",
			MetricsG: (*repo).GetGaugeList(),
			MetricsC: (*repo).GetCounterList(),
		}

		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html")

		tmpl, err := template.ParseFiles("./static/htmltemplates/main_page_template.html")

		if err != nil {
			fmt.Printf("error while loading template: %v\n", err)
			header = http.StatusInternalServerError
			w.WriteHeader(header)
		} else {

			w.WriteHeader(header)
			err = tmpl.Execute(w, data)

			if err != nil {
				fmt.Println("Error executing template:", err)
				//header = http.StatusInternalServerError
			}
		}

	})
}

// metric page
func GetPageM(repo *storage.Repositories) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg []byte
		header := http.StatusNotFound
		w.Header().Add("Content-Type", "text/plain")

		if chi.URLParam(r, "mName") != "" {

			switch strings.ToLower(chi.URLParam(r, "mType")) {
			case "gauge":
				if v, ok := (*repo).GetGauge(chi.URLParam(r, "mName")); ok {
					header = http.StatusOK
					msg = []byte(fmt.Sprint(v))
				}
			case "counter":
				if v, ok := (*repo).GetCounter(chi.URLParam(r, "mName")); ok {
					header = http.StatusOK
					msg = []byte(fmt.Sprint(v))
				}
			}

			w.WriteHeader(header)

			if len(msg) > 0 {
				if _, err := w.Write(msg); err != nil {
					log.Printf("Error while write msg: %s\n", err)
				}

			}
		}
	})
}

// processing all received data by "update" address
func UpdatePageM(repo *storage.Repositories) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// post message if needed
		var msg []byte
		// default header = badRequest
		header := http.StatusBadRequest

		// check method
		if r.Method == http.MethodPost {

			// check content type
			var rightContentT = true
			for i, v := range r.Header {
				if i == "Content-Type" && v[0] != "text/plain" {
					rightContentT = false
				}
			}

			if rightContentT {
				rcvdData := []string{
					chi.URLParam(r, "mType"),
					chi.URLParam(r, "mName"),
					chi.URLParam(r, "mValue"),
				}

				// processing received data
				err := rcvddataproc.UsePOSTData(rcvdData, repo)
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

		w.WriteHeader(header)
		if len(msg) > 0 {
			fmt.Printf("http msg: %s\n", msg)
		}
	})
}

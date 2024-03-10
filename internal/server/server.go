package server

import (
	//"net/http"
	"github.com/go-chi/chi/v5"
	hnd "github.com/ushanovsn/metrics/internal/handlers"
)


func ServerMux() *chi.Mux {
	r := chi.NewRouter()
	
	// route for post metrics
	r.Route("/update", func(r chi.Router){
		r.Post("/{mType}/{mName}/{mValue}", hnd.UpdatePageM)
	})
	
	// route for get
	r.Route("/", func(r chi.Router){
		r.Get("/", hnd.StartPage)
		r.Route("/value", func(r chi.Router){
			r.Get("/{mType}/{mName}", hnd.GetPageM)
		})
	})

	return r
}

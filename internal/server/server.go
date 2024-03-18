package server

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	hnd "github.com/ushanovsn/metrics/internal/handlers"
	"github.com/ushanovsn/metrics/internal/options"
)


func ServerRun() error {
	ServerOpt := options.ServerOptions{
		Net: options.NetAddress{
			Host: "",
			Port: 8080,
		},
	}

	InitFlag(&ServerOpt)
	InitEnv(&ServerOpt)
	

	return http.ListenAndServe(ServerOpt.Net.String(), ServerMux())
}

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

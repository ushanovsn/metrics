package server

import (
	"github.com/go-chi/chi/v5"
	hnd "github.com/ushanovsn/metrics/internal/handlers"
	"github.com/ushanovsn/metrics/internal/options"
	"github.com/ushanovsn/metrics/internal/storage"
	"net/http"
)

type serverData struct {
	opt   options.ServerOptions
	repo  storage.Repositories
	hndlr http.Handler
}

func ServerInit() *serverData {
	serverOpt := options.InitSrv()

	InitFlag(serverOpt)
	InitEnv(serverOpt)

	serverRepo := storage.Init()

	h := srvRouter(&serverRepo)

	return &serverData{
		opt:   *serverOpt,
		repo:  serverRepo,
		hndlr: h,
	}
}

func ServerRun(sd *serverData) error {
	return http.ListenAndServe(sd.opt.Net.String(), sd.hndlr)
}

func srvRouter(repo *storage.Repositories) *chi.Mux {

	r := chi.NewRouter()

	// route for post metrics
	r.Route("/update", func(r chi.Router) {
		r.Post("/{mType}/{mName}/{mValue}", hnd.UpdatePageM(repo))
	})

	// route for get
	r.Route("/", func(r chi.Router) {
		r.Get("/", hnd.StartPage(repo))
		r.Route("/value", func(r chi.Router) {
			r.Get("/{mType}/{mName}", hnd.GetPageM(repo))
		})
	})

	return r
}

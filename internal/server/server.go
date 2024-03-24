package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	hnd "github.com/ushanovsn/metrics/internal/handlers"
	"github.com/ushanovsn/metrics/internal/logger"
	"github.com/ushanovsn/metrics/internal/options"
	"github.com/ushanovsn/metrics/internal/storage"
	"net/http"
)

type serverData struct {
	opt   *options.ServerOptions
	repo  *storage.Repositories
	hndlr http.Handler
}

func ServerInit() *serverData {
	serverOpt := options.InitSrv()

	serverData := serverData{
		opt: serverOpt,
	}

	InitFlag(serverOpt)
	envErr := InitEnv(serverOpt)

	logger.InitLogger(&serverOpt.Logger)
	log := serverOpt.Logger.GetLogger()

	if envErr != nil {
		log.Errorf("Error when parsing environment var: %s\n", envErr)
	}

	serverRepo := storage.Init()

	serverData.repo = &serverRepo
	setRouter(&serverData)

	return &serverData
}

func ServerRun(sd *serverData) error {
	return http.ListenAndServe(sd.opt.Net.String(), sd.hndlr)
}

func setRouter(sd *serverData) {
	//repo *storage.Repositories,

	r := chi.NewRouter()

	// route for post metrics
	r.Route("/update", func(r chi.Router) {
		r.Post("/{mType}/{mName}/{mValue}", hnd.RespLogging(hnd.UpdatePageM(sd.repo), sd.opt.Logger.GetLogger()))
	})

	// route for get
	r.Route("/", func(r chi.Router) {
		r.Get("/", hnd.RespLogging(hnd.StartPage(sd.repo), sd.opt.Logger.GetLogger()))
		r.Route("/value", func(r chi.Router) {
			r.Get("/{mType}/{mName}", hnd.RespLogging(hnd.GetPageM(sd.repo), sd.opt.Logger.GetLogger()))
		})
	})

	sd.hndlr = r
}

func ServerStop(srv *serverData) {
	err := srv.opt.Logger.Stop()
	if err != nil {
		fmt.Printf("Error while stopping Agent: %s\n", err)
	}
}

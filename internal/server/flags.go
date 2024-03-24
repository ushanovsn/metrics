package server

import (
	"flag"
	"github.com/ushanovsn/metrics/internal/options"
)

func InitFlag(o *options.ServerOptions) {
	_ = flag.Value(&o.Net)

	flag.Var(&o.Net, "a", "Server net address host:port")
	flag.StringVar(&o.Logger.Level, "l", "info", "Logging level")

	flag.Parse()
}

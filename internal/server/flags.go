package server

import (
	"flag"
	"github.com/ushanovsn/metrics/internal/options"
)

func InitFlag(o *options.ServerOptions) {
	_ = flag.Value(&o.Net)

	flag.Var(&o.Net, "a", "Server net address host:port")

	flag.Parse()
}

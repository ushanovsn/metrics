package server

import (
	"flag"
	"github.com/ushanovsn/metrics/internal/options"
)


func InitFlag() {
	

	_ = flag.Value(&options.ServerOpt.Net)

	flag.Var(&options.ServerOpt.Net, "a", "Server net address host:port")

	flag.Parse()
}


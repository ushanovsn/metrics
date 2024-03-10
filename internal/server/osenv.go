package server

import (
	"github.com/ushanovsn/metrics/internal/options"
	"os"
)


func InitEnv() {
	if v, ok := os.LookupEnv("ADDRESS"); ok {
		options.ServerOpt.Net.Set(v)
	}
}
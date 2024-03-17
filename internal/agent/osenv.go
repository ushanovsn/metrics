package agent

import (
	"github.com/caarlos0/env/v6"
	"github.com/ushanovsn/metrics/internal/options"
	"os"
	//"log"
)


func InitEnv() {
	_ = env.Parse(&options.AgentOpt)

	if v, ok := os.LookupEnv("ADDRESS"); ok {
		options.AgentOpt.Net.Set(v)
	}
}
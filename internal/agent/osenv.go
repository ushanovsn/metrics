package agent

import (
	"github.com/caarlos0/env/v6"
	"github.com/ushanovsn/metrics/internal/options"
	"os"
	"log"
)



func InitEnv(o *options.AgentOptions) {
	_ = env.Parse(o)

	if v, ok := os.LookupEnv("ADDRESS"); ok {
		err := (*o).Net.Set(v)
		if err != nil {
			log.Printf("Error while set agent network address: %s\n", err.Error())
		}
	}
}
package agent

import (
	"github.com/caarlos0/env/v6"
	"github.com/ushanovsn/metrics/internal/options"
	"os"
)

func InitEnv(o *options.AgentOptions) error {
	_ = env.Parse(o)
	_ = env.Parse(o.Logger)

	if v, ok := os.LookupEnv("ADDRESS"); ok {
		err := (*o).Net.Set(v)
		if err != nil {
			return err
		}
	}
	return nil
}

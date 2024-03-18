package agent

import (
	"flag"
	"github.com/ushanovsn/metrics/internal/options"
	"log"
)

type flagOptions struct {
	reportInterval int `env:"REPORT_INTERVAL"`
	pollInterval   int `env:"POLL_INTERVAL"`
}

func InitFlag(o *options.AgentOptions) {
	flags := flagOptions{}
	_ = flag.Value(&o.Net)

	flag.Var(&o.Net, "a", "Server net address host:port")
	flag.IntVar(&flags.reportInterval, "r", 10, "Send metrics to server interval sec")
	flag.IntVar(&flags.pollInterval, "p", 2, "Update metrics interval sec")

	flag.Parse()

	if err := o.SetPolInt(flags.pollInterval); err != nil {
		log.Printf("Error when applying the Polling interval: %s\n", err.Error())
	}

	if err := o.SetRepInt(flags.reportInterval); err != nil {
		log.Printf("Error when applying the Polling interval: %s\n", err.Error())
	}

}

package agent

import (
	"flag"
	"github.com/ushanovsn/metrics/internal/options"
)

func InitFlag(o *options.AgentOptions) {
	//flags := flagOptions{}
	_ = flag.Value(&o.Net)

	flag.Var(&o.Net, "a", "Server net address host:port")
	//flag.IntVar(&flags.reportInterval, "r", 10, "Send metrics to server interval sec")
	flag.IntVar(&o.ReportInterval, "r", 10, "Send metrics to server interval sec")
	//flag.IntVar(&flags.pollInterval, "p", 2, "Update metrics interval sec")
	flag.IntVar(&o.PollInterval, "p", 2, "Update metrics interval sec")
	flag.StringVar(&o.Logger.Level, "l", "info", "Logging level")

	flag.Parse()
}

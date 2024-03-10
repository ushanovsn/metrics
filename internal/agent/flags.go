package agent

import (
	"flag"
	"github.com/ushanovsn/metrics/internal/options"
)



func InitFlag() {
	_ = flag.Value(&options.AgentOpt.Net)

	flag.Var(&options.AgentOpt.Net, "a", "Server net address host:port")
	flag.IntVar(&options.AgentOpt.ReportInterval, "r", 10, "Send metrics to server interval sec")
	flag.IntVar(&options.AgentOpt.PollInterval, "p", 2, "Update metrics interval sec")

	flag.Parse()
}


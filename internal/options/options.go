package options

import (
	"fmt"
	"strconv"
	"strings"
	"errors"
)
type NetAddress struct {
	Host string	
	Port int
}


type AgentOptions struct {
	Net NetAddress
	ReportInterval int		`env:"REPORT_INTERVAL"`
	PollInterval int		`env:"POLL_INTERVAL"`
}

var AgentOpt AgentOptions = AgentOptions{
	Net: NetAddress{
		Host: "localhost",
		Port: 8080,
	},
}




type ServerOptions struct {
	Net NetAddress
}

var ServerOpt ServerOptions = ServerOptions{
	Net: NetAddress{
		Host: "",
		Port: 8080,
	},
}




func (addr *NetAddress) String() string {
	return addr.Host + ":" + fmt.Sprint(addr.Port)
}

func (addr *NetAddress) Set(fVal string) error {
	vals := strings.Split(fVal, ":")
	if len(vals) != 2 {
		return errors.New("wrong flag values")
	}

	if v, err := strconv.Atoi(vals[1]); err == nil {
		addr.Host = vals[0]
		addr.Port = v
		return nil
	} else {
		return errors.New("wrong port value")
	}
}
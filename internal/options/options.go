package options

import (
	"fmt"
	"strconv"
	"strings"
	"errors"
)

type NetAddress struct {
	host string	
	port int
}


type agentOptions struct {
	Net NetAddress
	reportInterval int		`env:"REPORT_INTERVAL"`
	pollInterval int		`env:"POLL_INTERVAL"`
}

var AgentOpt agentOptions = agentOptions{
	Net: NetAddress{
		host: "localhost",
		port: 8080,
	},
	reportInterval: 10,
	pollInterval: 2,
}




type serverOptions struct {
	Net NetAddress
}

var ServerOpt serverOptions = serverOptions{
	Net: NetAddress{
		host: "",
		port: 8080,
	},
}



func (n *NetAddress)String() string {
	return n.host + ":" + fmt.Sprint(n.port)
}

func (n *NetAddress)Set(s string) error {
	vals := strings.Split(s, ":")
	if len(vals) != 2 {
		return errors.New("wrong addres string")
	}

	if v, err := strconv.Atoi(vals[1]); err == nil {
		n.host = vals[0]
		n.port = v
	} else {
		return errors.New("wrong port value")
	}

	return nil
}


func (a *agentOptions)SetPolInt(v int) error {
	if v > 0 {
		a.pollInterval = v
	} else {
		return errors.New("wrong value of poll interval")
	}
	return nil
}

func (a *agentOptions)SetRepInt(v int) error {
	if v > 0 {
		a.reportInterval = v
	} else {
		return errors.New("wrong value of report interval")
	}
	return nil
}

func (a *agentOptions)GetPollInt() int {
	return a.pollInterval
}

func (a *agentOptions)GetRepInt() int {
	return a.reportInterval
}
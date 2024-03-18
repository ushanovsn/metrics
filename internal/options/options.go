package options

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

type AgentOptions struct {
	Net            NetAddress
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

type ServerOptions struct {
	Net NetAddress
}

func (n *NetAddress) String() string {
	return n.Host + ":" + fmt.Sprint(n.Port)
}

func (n *NetAddress) Set(s string) error {
	vals := strings.Split(s, ":")
	if len(vals) != 2 {
		return errors.New("wrong addres string")
	}

	if v, err := strconv.Atoi(vals[1]); err == nil {
		n.Host = vals[0]
		n.Port = v
	} else {
		return errors.New("wrong port value")
	}

	return nil
}

func (a *AgentOptions) SetPolInt(v int) error {
	if v > 0 {
		a.PollInterval = v
	} else {
		return errors.New("wrong value of poll interval")
	}
	return nil
}

func (a *AgentOptions) SetRepInt(v int) error {
	if v > 0 {
		a.ReportInterval = v
	} else {
		return errors.New("wrong value of report interval")
	}
	return nil
}

func (a *AgentOptions) GetPollInt() int {
	return a.PollInterval
}

func (a *AgentOptions) GetRepInt() int {
	return a.ReportInterval
}

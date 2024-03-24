package options

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

type LoggerOpt struct {
	FileName string `env:"AGENT_LOG_FILE_NAME"`
	Level    string `env:"LOG_LEVEL"`
	Report   bool   `env:"LOG_REPORT"`
	file     *os.File
	logger   *logrus.Logger
}

type AgentOptions struct {
	Net            NetAddress
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
	Logger         LoggerOpt
}

type ServerOptions struct {
	Net    NetAddress
	Logger LoggerOpt
}

func initNetParam() *NetAddress {
	return &NetAddress{
		Host: "localhost",
		Port: 8080,
	}
}

func InitAg() *AgentOptions {
	return &AgentOptions{
		Net:            *initNetParam(),
		ReportInterval: 10,
		PollInterval:   2,
		Logger: LoggerOpt{
			FileName: "log_agent",
			Level:    "debug",
			Report:   true,
		},
	}
}

func InitSrv() *ServerOptions {
	return &ServerOptions{
		Net: *initNetParam(),
		Logger: LoggerOpt{
			FileName: "log_server",
			Level:    "debug",
			Report:   true,
		},
	}
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

func (log *LoggerOpt) GetLogger() *logrus.Logger {
	return log.logger
}

func (log *LoggerOpt) Stop() error {
	if log.file != nil {
		return log.file.Close()
	}
	return nil
}

func (log *LoggerOpt) SetFile(f *os.File) {
	log.file = f
}

func (log *LoggerOpt) SetLogger(l *logrus.Logger) {
	log.logger = l
}

package server

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

type FlagParameters struct {
	Net NetAddress
}

var FlagParam FlagParameters = FlagParameters{
	Net: NetAddress{
		Host: "",
		Port: 8080,
	},
}

func FlagInit() {
	_ = flag.Value(&FlagParam.Net)

	flag.Var(&FlagParam.Net, "a", "Server net address host:port")

	flag.Parse()
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

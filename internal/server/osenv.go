package server

import (
	"github.com/ushanovsn/metrics/internal/options"
	"os"
	"log"
)


func InitEnv() {
	if v, ok := os.LookupEnv("ADDRESS"); ok {
		err := options.ServerOpt.Net.Set(v)
		if err != nil {
			log.Printf("Error while set server network address: %s\n", err.Error())
		}
	}
}
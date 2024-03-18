package server

import (
	"github.com/ushanovsn/metrics/internal/options"
	"os"
	"log"
)


func InitEnv(o *options.ServerOptions) {
	if v, ok := os.LookupEnv("ADDRESS"); ok {
		err := (*o).Net.Set(v)
		if err != nil {
			log.Printf("Error while set server network address: %s\n", err.Error())
		}
	}
}
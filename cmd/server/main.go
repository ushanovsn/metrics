package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/server"
)


func main() {
	fmt.Printf("Server starting...\n")

	server.FlagInit()

    if err := server.ServerRun(); err != nil {
        panic(err)
    }
	fmt.Printf("Server stopped! \n")
}

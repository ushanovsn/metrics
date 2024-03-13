package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/server"
)


func main() {
	fmt.Printf("Server starting...\n")

	server.InitFlag()
	server.InitEnv()

	if err := server.ServerRun(); err != nil {
		fmt.Printf("Server stopping with error: %s \n", err)
	}

	fmt.Printf("Server stopped! \n")
}

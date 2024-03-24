package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/server"
)

func main() {
	fmt.Printf("Server init...\n")

	srv := server.ServerInit()

	fmt.Printf("Server starting...\n")

	if err := server.ServerRun(srv); err != nil {
		fmt.Printf("Server stopping with error: %s\n", err)
	}

	server.ServerStop(srv)

	fmt.Printf("Server stopped!\n")
}

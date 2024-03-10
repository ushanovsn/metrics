package main

import (
	"fmt"
	"net/http"
	"github.com/ushanovsn/metrics/internal/server"
)


func main() {
	fmt.Printf("Server starting...\n")

	err := http.ListenAndServe(":8080", server.ServerMux())

	if err != nil {
		panic(err)
	}

	fmt.Printf("Server stopped! \n")
}

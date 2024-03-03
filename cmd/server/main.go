package main

import (
	"fmt"
	"net/http"
	hnd "github.com/ushanovsn/metrics/internal/handlers"
)


func main() {
	fmt.Printf("Server starting...\n")

	err := http.ListenAndServe(":8080", hnd.ServerMux())

	if err != nil {
		panic(err)
	}

	fmt.Printf("Server stopped! \n")
}

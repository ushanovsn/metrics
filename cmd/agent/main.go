package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/agent"
)


func main() {
	fmt.Printf("Starting client \n")

	agent.StartAgent()
}

package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/agent"
)


func main() {
	fmt.Printf("Starting client \n")

	if err := agent.AgentRun(); err != nil {
		fmt.Printf("Agent stopping with error: %s\n", err)
	}

	fmt.Printf("Agent stopped!\n")
}

package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/agent"
)

func main() {
	fmt.Printf("Agent init...\n")

	agnt := agent.AgentInit()

	fmt.Printf("Agent starting...\n")

	if err := agent.AgentRun(agnt); err != nil {
		fmt.Printf("Agent stopping with error: %s\n", err)
	}

	fmt.Printf("Agent stopped!\n")
}

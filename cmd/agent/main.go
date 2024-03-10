package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/agent"
)


func main() {
	fmt.Printf("Starting client \n")

	agent.FlagInit()

    if err := agent.AgentRun(); err != nil {
        panic(err)
    }
	
	fmt.Printf("Agent stopped! \n")
}

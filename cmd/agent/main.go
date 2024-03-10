package main

import (
	"fmt"
	"github.com/ushanovsn/metrics/internal/agent"
)


func main() {
	fmt.Printf("Starting client \n")

	agent.InitFlag()
	agent.InitEnv()

    if err := agent.AgentRun(); err != nil {
        panic(err)
    }
	
	fmt.Printf("Agent stopped! \n")
}

package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent/remoteagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/server/restapi/services"
)

func main() {
	ctx := context.Background()

	a2aServerAddress := startWeatherAgentServer()

	remoteAgent, err := remoteagent.New(remoteagent.A2AConfig{
		Name:            "A2A Demo Agent",
		AgentCardSource: a2aServerAddress,
	})
	if err != nil {
		log.Fatalf("Failed to create a remote agent: %v", err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(remoteAgent),
	}

	l := full.NewLauncher()
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

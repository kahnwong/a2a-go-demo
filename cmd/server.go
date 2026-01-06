package cmd

import (
	"context"
	"log"
	"strconv"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/session"
)

func startA2AServer(ctx context.Context, port int, agentInstance agent.Agent, serverName string) {
	webLauncher := web.NewLauncher(a2a.NewLauncher())
	_, err := webLauncher.Parse([]string{
		"--port", strconv.Itoa(port),
		"a2a", "--a2a_agent_url", "http://localhost:" + strconv.Itoa(port),
	})
	if err != nil {
		log.Fatalf("launcher.Parse() error = %v", err)
	}

	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(agentInstance),
		SessionService: session.InMemoryService(),
	}

	log.Printf("Starting A2A %s server on port %d\n", serverName, port)
	if err := webLauncher.Run(ctx, config); err != nil {
		log.Fatalf("webLauncher.Run() error = %v", err)
	}
}

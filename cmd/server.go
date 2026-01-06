package cmd

import (
	"context"
	"strconv"

	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).Msg("launcher.Parse() failed")
	}

	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(agentInstance),
		SessionService: session.InMemoryService(),
	}

	log.Info().Str("server", serverName).Int("port", port).Msg("Starting A2A server")
	if err := webLauncher.Run(ctx, config); err != nil {
		log.Fatal().Err(err).Msg("webLauncher.Run() failed")
	}
}

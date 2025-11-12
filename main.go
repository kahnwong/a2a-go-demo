package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/remoteagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func newWeatherAgent() (agent.Agent, error) {
	remoteAgent, err := remoteagent.New(remoteagent.A2AConfig{
		Name:            "weather_agent",
		Description:     "Agent that checks current weather in a given city.",
		AgentCardSource: "http://localhost:8001",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create remote weather agent: %w", err)
	}
	return remoteAgent, nil
}

func newTimeAgent() (agent.Agent, error) {
	remoteAgent, err := remoteagent.New(remoteagent.A2AConfig{
		Name:            "time_agent",
		Description:     "Agent that checks current time in a given city.",
		AgentCardSource: "http://localhost:8002",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create remote time agent: %w", err)
	}
	return remoteAgent, nil
}

func newRootAgent(ctx context.Context, weatherAgent, timeAgent agent.Agent) (agent.Agent, error) {
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{})
	if err != nil {
		return nil, err
	}
	return llmagent.New(llmagent.Config{
		Name:        "root_agent",
		Model:       model,
		Description: "Agent to answer questions about the time and weather in a city.",
		Instruction: "I can answer your questions about the time and weather in a city.",
		SubAgents:   []agent.Agent{weatherAgent, timeAgent},
		Tools:       []tool.Tool{},
	})
}

func main() {
	ctx := context.Background()

	// init remote agents
	weatherAgent, err := newWeatherAgent()
	if err != nil {
		log.Fatalf("Failed to create weather agent: %v", err)
	}

	timeAgent, err := newTimeAgent()
	if err != nil {
		log.Fatalf("Failed to create time agent: %v", err)
	}

	// init root agent
	rootAgent, err := newRootAgent(ctx, weatherAgent, timeAgent)
	if err != nil {
		log.Fatalf("Failed to create root agent: %v", err)
	}

	// start server
	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(rootAgent),
	}

	l := full.NewLauncher()
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

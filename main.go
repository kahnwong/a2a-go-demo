package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/remoteagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// init root agent
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	rootAgent, err := llmagent.New(llmagent.Config{
		Name:        "weather_time_agent",
		Model:       model,
		Description: "Agent to answer questions about the time and weather in a city.",
		Instruction: "I can answer your questions about the time and weather in a city.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// init remote agents
	weatherServerAddress := startWeatherAgentServer()
	weatherAgent, err := remoteagent.New(remoteagent.A2AConfig{
		Name:            "A2A Weather Agent",
		AgentCardSource: weatherServerAddress,
	})
	if err != nil {
		log.Fatalf("Failed to create a weather agent: %v", err)
	}

	timeServerAddress := startTimeAgentServer()
	timeAgent, err := remoteagent.New(remoteagent.A2AConfig{
		Name:            "A2A Time Agent",
		AgentCardSource: timeServerAddress,
	})
	if err != nil {
		log.Fatalf("Failed to create a weather agent: %v", err)
	}

	// start server
	agentLoader, err := services.NewMultiAgentLoader(
		rootAgent,
		weatherAgent,
		timeAgent,
	)
	if err != nil {
		log.Fatalf("Failed to create agent loader: %v", err)
	}

	config := &adk.Config{
		AgentLoader: agentLoader,
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

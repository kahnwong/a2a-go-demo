package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

type Input struct {
	City string `json:"city" jsonschema:"city name"`
}

type Output struct {
	WeatherSummary string `json:"weather_summary" jsonschema:"weather summary in the given city"`
}

func GetWeather(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	return nil, Output{
		WeatherSummary: fmt.Sprintf("Today in %q is sunny\n", input.City),
	}, nil
}

func localMCPTransport(ctx context.Context) mcp.Transport {
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Run in-memory MCP server.
	server := mcp.NewServer(&mcp.Implementation{Name: "weather_server", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "get_weather", Description: "returns weather in the given city"}, GetWeather)
	_, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	return clientTransport
}

func main() {
	ctx := context.Background()

	// init model
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// init tools
	transport := localMCPTransport(ctx)

	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
	}

	// init agent
	agent, err := llmagent.New(llmagent.Config{
		Name:        "agent_a",
		Model:       model,
		Description: "Helper agent.",
		Instruction: "You are a helpful assistant that helps users with various tasks.",
		Toolsets: []tool.Toolset{
			mcpToolSet,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// start server
	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}

	l := full.NewLauncher()
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

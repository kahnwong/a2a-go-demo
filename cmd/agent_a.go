package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
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

	server := mcp.NewServer(&mcp.Implementation{Name: "weather_server", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "get_weather", Description: "returns weather in the given city"}, GetWeather)
	_, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect MCP server")
	}

	return clientTransport
}

func WeatherAgent(ctx context.Context) agent.Agent {
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create model")
	}

	transport := localMCPTransport(ctx)

	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create MCP tool set")
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "agent_a",
		Model:       model,
		Description: "Agent to answer questions about the weather in a city.",
		Instruction: "You are a helpful assistant that helps users with various tasks.",
		Toolsets: []tool.Toolset{
			mcpToolSet,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create agent")
	}

	return agent
}

var agentACmd = &cobra.Command{
	Use:   "agent-a",
	Short: "Start agent A (weather agent)",
	Long:  `Starts the weather agent on port 8001`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		startA2AServer(ctx, 8001, WeatherAgent(ctx), "weather agent")
	},
}

func init() {
	rootCmd.AddCommand(agentACmd)
}

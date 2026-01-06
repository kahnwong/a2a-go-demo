package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
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
		log.Fatal(err)
	}

	return clientTransport
}

func WeatherAgent(ctx context.Context) agent.Agent {
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create a model: %v", err)
	}

	transport := localMCPTransport(ctx)

	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
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
		log.Fatalf("Failed to create agent: %v", err)
	}

	return agent
}

var agentACmd = &cobra.Command{
	Use:   "agent-a",
	Short: "Start agent A (weather agent)",
	Long:  `Starts the weather agent on port 8001`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		port := 8001
		webLauncher := web.NewLauncher(a2a.NewLauncher())
		_, err := webLauncher.Parse([]string{
			"--port", strconv.Itoa(port),
			"a2a", "--a2a_agent_url", "http://localhost:" + strconv.Itoa(port),
		})
		if err != nil {
			log.Fatalf("launcher.Parse() error = %v", err)
		}

		config := &launcher.Config{
			AgentLoader:    agent.NewSingleLoader(WeatherAgent(ctx)),
			SessionService: session.InMemoryService(),
		}

		log.Printf("Starting A2A weather agent server on port %d\n", port)
		if err := webLauncher.Run(context.Background(), config); err != nil {
			log.Fatalf("webLauncher.Run() error = %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentACmd)
}

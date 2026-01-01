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
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

type TimeInput struct {
	City string `json:"city" jsonschema:"city name"`
}

type TimeOutput struct {
	TimeSummary string `json:"time_summary" jsonschema:"time summary in the given city"`
}

func GetTime(ctx context.Context, req *mcp.CallToolRequest, input TimeInput) (*mcp.CallToolResult, TimeOutput, error) {
	return nil, TimeOutput{
		TimeSummary: fmt.Sprintf("Current time in %q is 5:00 AM\n", input.City),
	}, nil
}

func localAgentBMCPTransport(ctx context.Context) mcp.Transport {
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{Name: "time_server", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "get_Time", Description: "returns time in the given city"}, GetTime)
	_, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	return clientTransport
}

func TimeAgent(ctx context.Context) agent.Agent {
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create a model: %v", err)
	}

	transport := localAgentBMCPTransport(ctx)

	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "agent_b",
		Model:       model,
		Description: "Agent to answer questions about the time in a city.",
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

var agentBCmd = &cobra.Command{
	Use:   "agent-b",
	Short: "Start agent B (time agent)",
	Long:  `Starts the time agent on port 8002`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		port := 8002
		launcher := web.NewLauncher(a2a.NewLauncher())
		_, err := launcher.Parse([]string{
			"--port", strconv.Itoa(port),
			"a2a", "--a2a_agent_url", "http://localhost:" + strconv.Itoa(port),
		})
		if err != nil {
			log.Fatalf("launcher.Parse() error = %v", err)
		}

		config := &adk.Config{
			AgentLoader:    services.NewSingleAgentLoader(TimeAgent(ctx)),
			SessionService: session.InMemoryService(),
		}

		log.Printf("Starting A2A time agent server on port %d\n", port)
		if err := launcher.Run(context.Background(), config); err != nil {
			log.Fatalf("launcher.Run() error = %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentBCmd)
}

package agents

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
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
		log.Fatal().Err(err).Msg("Failed to connect MCP server")
	}

	return clientTransport
}

func TimeAgent(ctx context.Context) agent.Agent {
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create model")
	}

	transport := localAgentBMCPTransport(ctx)

	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create MCP tool set")
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
		log.Fatal().Err(err).Msg("Failed to create agent")
	}

	return agent
}

func StartAgentB() {
	ctx := context.Background()

	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown tracer")
		}
	}()

	startA2AServer(ctx, 8002, TimeAgent(ctx), "time agent")
}

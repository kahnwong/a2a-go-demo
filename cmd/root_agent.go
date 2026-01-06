package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/remoteagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func newWeatherAgent() (agent.Agent, error) {
	remoteAgent, err := remoteagent.NewA2A(remoteagent.A2AConfig{
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
	remoteAgent, err := remoteagent.NewA2A(remoteagent.A2AConfig{
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

var rootAgentCmd = &cobra.Command{
	Use:   "root-agent",
	Short: "Start root agent",
	Long:  `Starts the root agent that coordinates weather and time agents`,
	Run: func(cmd *cobra.Command, args []string) {
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

		weatherAgent, err := newWeatherAgent()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create weather agent")
		}

		timeAgent, err := newTimeAgent()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create time agent")
		}

		rootAgent, err := newRootAgent(ctx, weatherAgent, timeAgent)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create root agent")
		}

		config := &launcher.Config{
			AgentLoader: agent.NewSingleLoader(rootAgent),
		}

		l := full.NewLauncher()
		launchArgs := []string{"web", "api", "webui"}
		if len(args) > 0 {
			launchArgs = args
		}
		err = l.Execute(ctx, config, launchArgs)
		if err != nil {
			log.Fatal().Err(err).Str("syntax", l.CommandLineSyntax()).Msg("Run failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(rootAgentCmd)
}

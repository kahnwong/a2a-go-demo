package cmd

import (
	"github.com/kahnwong/a2a-demo/agents"
	"github.com/spf13/cobra"
)

var agentBCmd = &cobra.Command{
	Use:   "agent-b",
	Short: "Start agent B (time agent)",
	Long:  `Starts the time agent on port 8002`,
	Run: func(cmd *cobra.Command, args []string) {
		agents.StartAgentB()
	},
}

func init() {
	rootCmd.AddCommand(agentBCmd)
}

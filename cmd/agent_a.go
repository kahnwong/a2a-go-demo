package cmd

import (
	"github.com/kahnwong/a2a-demo/agents"
	"github.com/spf13/cobra"
)

var agentACmd = &cobra.Command{
	Use:   "agent-a",
	Short: "Start agent A (weather agent)",
	Long:  `Starts the weather agent on port 8001`,
	Run: func(cmd *cobra.Command, args []string) {
		agents.StartAgentA()
	},
}

func init() {
	rootCmd.AddCommand(agentACmd)
}

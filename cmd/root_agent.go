package cmd

import (
	"github.com/kahnwong/a2a-demo/agents"
	"github.com/spf13/cobra"
)

var rootAgentCmd = &cobra.Command{
	Use:   "root-agent",
	Short: "Start root agent",
	Long:  `Starts the root agent that coordinates weather and time agents`,
	Run: func(cmd *cobra.Command, args []string) {
		agents.StartRootAgent(args)
	},
}

func init() {
	rootCmd.AddCommand(rootAgentCmd)
}

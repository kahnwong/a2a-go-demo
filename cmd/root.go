package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "a2a-go-demo",
	Short: "A2A Go Demo - Agent communication demo",
	Long:  `A2A Go Demo is an application that demonstrates agent-to-agent communication using Google ADK.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

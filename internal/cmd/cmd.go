package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Short: "Campaign service",
}

func Execute() {
	RootCmd.AddCommand(versionCmd)

	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(migrationCmd)

	if err := RootCmd.Execute(); err != nil {
		log.Printf("failed to start program, err: %v", err)
		os.Exit(-1)
	}
}

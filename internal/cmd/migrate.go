package cmd

import (
	"campaign/internal/db"

	"github.com/spf13/cobra"
)

var (
	migrationConfigPath string
)

func init() {
	migrationCmd.PersistentFlags().StringVarP(&migrationConfigPath, "config", "c", "config.yml", "config to file path")
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Long:  `Migrate database`,
	Run: func(cmd *cobra.Command, args []string) {
		username := "user"
		password := "password"
		dbname := "campaign"
		host := "192.168.50.12:5433"

		dbClient := db.Get(host, dbname, username, password)

		db.RunMigrations("scripts", dbClient)
	},
}

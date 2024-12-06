package cmd

import (
	"campaign/internal/db"
	"os"
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

		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "192.168.50.12:5433"
		}

		dbClient := db.Get(dbHost, dbname, username, password)

		db.RunMigrations("scripts", dbClient)
	},
}

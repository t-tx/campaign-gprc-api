package cmd

import (
	"campaign/internal/cache"
	"campaign/internal/db"
	"campaign/internal/db/repositories"
	"campaign/internal/handler"
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"os"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start campaign application",
	Long:  `Start campaign application`,
	Run: func(cmd *cobra.Command, args []string) {
		jwtSecret := []byte("jwtSecret")
		server := pkg.NewGrpcServer(func(s *grpc.Server) {
			username := "user"
			password := "password"
			dbname := "campaign"

			dbHost := os.Getenv("DB_HOST")
			if dbHost == "" {
				dbHost = "192.168.50.12:5433"
			}
			dbClient := db.Get(dbHost, dbname, username, password)
			repository := repositories.NewRepository(dbClient)

			cacheHost := os.Getenv("CACHE_HOST")
			if cacheHost == "" {
				cacheHost = "192.168.50.12:6379"
			}
			cacheClient := cache.NewClient(cacheHost)

			servicepb.RegisterCampaignServiceServer(s, handler.New(repository, cacheClient, jwtSecret))
		}, "campaign", 8000)

		server.AddInterceptors(pkg.AuthInterceptor("token", jwtSecret))
		server.AddInterceptors(pkg.LogInterceptor())

		server.RegisterGW(servicepb.RegisterCampaignServiceHandlerFromEndpoint)

		server.Run()
	},
}

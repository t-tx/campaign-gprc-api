package main

import (
	"campaign/internal/cmd"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	revision string
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	viper.SetDefault("author", "trungtvq.work@gmail.com")

	cmd.SetRevision(revision)
	cmd.Execute()
}

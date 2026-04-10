package main

import (
	"context"
	"log"

	"github.com/Danil-Ivonin/GrintexRates/internal/config"
	"github.com/Danil-Ivonin/GrintexRates/internal/http/client"
	"github.com/Danil-Ivonin/GrintexRates/internal/observability"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, err := observability.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync() //nolint:errcheck

	err = config.Load()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}
	url := viper.GetString("grinex.url")
	timeout := viper.GetDuration("grinex.tumeout")
	cl := client.New(url, timeout)
	fetch, err := cl.Fetch(context.Background())
	logger.Info("fetch finished", zap.Any("fetch", fetch))
}

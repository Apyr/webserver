package main

import (
	"os"
	"time"
	"webserver/config"
	"webserver/server"

	"go.uber.org/zap"
)

func loadConfig() (*config.Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	err = cfg.LoadFromYAML(data)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func newLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	logger.Core().Enabled(zap.DebugLevel)
	if err != nil {
		panic(err.Error())
	}
	return logger.Sugar()
}

const timeout = 10 * time.Second

func main() {
	logger := newLogger()
	defer logger.Sync()

	config, err := loadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	servers := server.NewServers(config, logger)

	logger.Infoln("starting servers")
	stop, errors := server.StartServers(servers, logger, timeout)

	go server.Wait(stop)

	for err := range errors {
		if err != nil {
			logger.Error(err)
		}
	}
}

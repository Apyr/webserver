package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"webserver/config"
	"webserver/server"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func newLogger(config *config.Config) *zap.SugaredLogger {
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.TimeKey = "time"
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logCfg.Level.SetLevel(config.GetLogLevel())
	if config.LogFile != "/dev/null" {
		logCfg.OutputPaths = []string{"stdout", config.LogFile}
	} else {
		logCfg.OutputPaths = []string{"stdout"}
	}
	logCfg.ErrorOutputPaths = logCfg.OutputPaths

	logger, err := logCfg.Build()
	if err != nil {
		panic(err.Error())
	}
	return logger.Sugar()
}

func watchForFile(fileName string) (bool, error) {
	// fs watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return false, err
	}
	defer watcher.Close()

	err = watcher.Add(fileName)
	if err != nil {
		return false, err
	}

	// Ctrl-C signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	// wait for any event
	select {
	case <-signals:
		return false, nil
	case <-watcher.Events:
		return true, nil
	case err := <-watcher.Errors:
		return false, err
	}
}

func main() {
	continueExec := true
	for continueExec {
		config, err := loadConfig()
		if err != nil {
			panic(err.Error())
		}

		logger := newLogger(config)
		defer logger.Sync()

		servers := server.NewServers(config, logger)
		servers.Start()

		continueExec, err = watchForFile("config.yaml")
		if err != nil {
			logger.Error(err)
		} else {
			if continueExec {
				logger.Info("config.yaml changed, reloading...")
			} else {
				logger.Info("interrupt signal received, stoping...")
			}
		}

		servers.Stop()
		for err := range servers.Errors {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info(err)
			} else {
				logger.Error(err)
			}
		}
	}
}

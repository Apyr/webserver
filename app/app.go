package app

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"webserver/config"
	"webserver/server"

	"github.com/fsnotify/fsnotify"
)

type App struct {
	logger     *slog.Logger
	config     *config.Config
	watcher    *fsnotify.Watcher
	interrupts chan os.Signal
}

func NewApp() *App {
	app := &App{
		logger:     slog.Default(),
		interrupts: make(chan os.Signal, 1),
	}
	signal.Notify(app.interrupts, os.Interrupt)
	if err := app.newWatcher(); err != nil {
		app.logger.Error("watcher creation error", slog.Any("error", err))
	}
	return app
}

func (app *App) Close() {
	if app.watcher != nil {
		app.watcher.Close()
	}
	signal.Stop(app.interrupts)
	close(app.interrupts)
}

func (app *App) Run() {
	for app.mainProcess() {
	}
}

func (app *App) newWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Add("config.yaml"); err != nil {
		return err
	}
	app.watcher = watcher
	return nil
}

func (app *App) mainProcess() bool {
	if err := app.loadConfig(); err != nil {
		app.logger.Error("config loading error", slog.Any("error", err))
		return app.waitForConfig()
	}

	if app.config != nil {
		closeLog := app.newLogger()
		defer closeLog()
	}

	servers := server.NewServers(app.config, app.logger)
	servers.Start()
	defer servers.Stop()

	for {
		select {
		case err := <-app.watcher.Errors:
			app.logger.Error("watching error", slog.Any("error", err))
			return true
		case <-app.watcher.Events:
			app.logger.Info("config.yaml changed, reloading...")
			return true
		case <-app.interrupts:
			app.logger.Info("interrupted")
			return false
		case err := <-servers.Errors:
			if err == nil || errors.Is(err, http.ErrServerClosed) {
				app.logger.Info("server closed")
			} else {
				app.logger.Error("server error", slog.Any("error", err))
			}
		}
	}
}

func (app *App) waitForConfig() bool {
	select {
	case err := <-app.watcher.Errors:
		app.logger.Error("watching error", slog.Any("error", err))
		return true
	case <-app.watcher.Events:
		app.logger.Info("config.yaml changed, reloading...")
		return true
	case <-app.interrupts:
		app.logger.Info("interrupted")
		return false
	}
}

func (app *App) loadConfig() error {
	app.config = nil
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	var cfg config.Config
	if err = cfg.LoadFromYAML(data); err != nil {
		return err
	}
	app.config = &cfg
	return nil
}

func (app *App) newLogger() func() {
	var out io.Writer = os.Stdout
	close := func() {}

	if app.config.LogFile != "/dev/null" {
		file, err := os.OpenFile(app.config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			app.logger.Error("failed to open log file", slog.Any("error", err))
		} else {
			out = io.MultiWriter(os.Stdout, file)
			close = func() {
				file.Close()
			}
		}
	}

	app.logger = slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: app.config.GetLogLevel()}))
	return close
}

package server

import (
	"log/slog"
	"net/http"

	"github.com/apyr/webserver/config"
	"github.com/apyr/webserver/handlers"
	httplogger "github.com/shogo82148/go-http-logger"
)

func makeHttpHandler(cfg *config.Config, logger *slog.Logger, isHTTPS bool) http.Handler {
	handler := handlers.NewEndpointsHandler(cfg.Endpoints, isHTTPS)
	requestLogger := httplogger.NewSlogLogger(slog.LevelDebug, "request", logger)
	return httplogger.LoggingHandler(requestLogger, handler)
}

package server

import (
	"net/http"
	"webserver/config"
	"webserver/handlers"

	httplogger "github.com/shogo82148/go-http-logger"
	"go.uber.org/zap"
)

func makeHttpHandler(cfg *config.Config, logger *zap.SugaredLogger, isHTTPS bool) http.Handler {
	handler := handlers.NewEndpointsHandler(cfg.Endpoints, isHTTPS)
	return wrapInLogger(handler, logger)
}

func wrapInLogger(handler http.Handler, logger *zap.SugaredLogger) http.Handler {
	return httplogger.LoggingHandler(httplogger.LoggerFunc(func(l httplogger.Attrs, r *http.Request) {
		logger.Debugw("request",
			"method", r.Method,
			"path", r.RequestURI,
			"host", r.Host,
			"requestSize", l.RequestSize(),
			"responseStatus", l.Status(),
			"responseSize", l.ResponseSize(),
			"responseContentType", l.Header().Get("Content-Type"),
		)
	}), handler)
}

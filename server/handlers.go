package server

import (
	"net/http"
	"webserver/config"
	"webserver/handlers"

	httplogger "github.com/shogo82148/go-http-logger"
	"go.uber.org/zap"
)

func makeHttpHandler(cfg *config.Config, logger *zap.SugaredLogger, isHTTPS bool) http.Handler {
	router := handlers.NewRouter()
	for _, endpoint := range cfg.Endpoints {
		if endpoint.Disabled {
			continue
		}
		handler := handlers.NewEndpointHandler(endpoint, isHTTPS)
		router.Add(endpoint.URL.URL, handler)
	}
	return wrapInLogger(router, logger)
}

func wrapInLogger(handler http.Handler, logger *zap.SugaredLogger) http.Handler {
	return httplogger.LoggingHandler(httplogger.LoggerFunc(func(l httplogger.Attrs, r *http.Request) {
		logger.Debugw("request",
			"requestURI", r.RequestURI,
			"responseSize", l.ResponseSize(),
			"status", l.Status(),
			"method", r.Method,
			"contentType", l.Header().Get("Content-Type"),
		)
	}), handler)
}

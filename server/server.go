package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/apyr/webserver/config"
	"github.com/apyr/webserver/tlsutils"
)

type ServerError struct {
	Server *http.Server
	err    error
}

func (err ServerError) Unwrap() error {
	return err.err
}

func (err ServerError) Error() string {
	return fmt.Sprintf("error at server %s: %s", err.Server.Addr, err.err.Error())
}

const defaultTimeout = 10 * time.Second

type Servers struct {
	Logger *slog.Logger
	Config *config.Config

	ShutdownTimeout time.Duration
	Errors          chan error

	HttpServer  *http.Server
	HttpsServer *http.Server

	wg *sync.WaitGroup
}

type logWriter struct {
	logger *slog.Logger
}

func (writer logWriter) Write(p []byte) (int, error) {
	writer.logger.Error(string(p))
	return len(p), nil
}

func NewServers(config *config.Config, logger *slog.Logger) Servers {
	useHTTPS := config.HTTPS()
	requestsLogger := logger.WithGroup("requests")
	httpHandler := makeHttpHandler(config, requestsLogger, false)
	httpsHandler := makeHttpHandler(config, requestsLogger, true)

	servers := Servers{
		Logger:          logger,
		Config:          config,
		ShutdownTimeout: defaultTimeout,
		Errors:          make(chan error),
		wg:              &sync.WaitGroup{},
	}

	errLogger := log.New(logWriter{logger}, "", 0)

	if useHTTPS {
		var tlsCfg *tls.Config
		tlsCfg, httpHandler = tlsutils.NewTLSConfig(
			jsonFileCache(config.CertsFile),
			httpHandler,
			config.AcmeHosts(),
			config.SelfHosts(),
		)

		servers.HttpsServer = &http.Server{
			Addr:      fmt.Sprintf("0.0.0.0:%d", config.Ports.HTTPS),
			Handler:   httpsHandler,
			TLSConfig: tlsCfg,
			ErrorLog:  errLogger,
		}
	}

	servers.HttpServer = &http.Server{
		Addr:     fmt.Sprintf("0.0.0.0:%d", config.Ports.HTTP),
		Handler:  httpHandler,
		ErrorLog: errLogger,
	}

	return servers
}

// asynchronously starting servers
func (servers Servers) Start() {
	servers.Logger.Info("starting servers")

	go servers.start(servers.HttpServer)
	if servers.HttpsServer != nil {
		go servers.start(servers.HttpsServer)
	}
}

func (servers Servers) start(server *http.Server) {
	servers.wg.Add(1)
	defer servers.wg.Done()

	var err error
	if server.TLSConfig != nil {
		err = server.ListenAndServeTLS("", "")
	} else {
		err = server.ListenAndServe()
	}

	err = ServerError{server, err}
	servers.Errors <- err
}

// asynchronously stoping servers
func (servers Servers) Stop() {
	servers.Logger.Info("stoping servers")

	go func() {
		servers.wg.Wait()
		close(servers.Errors)
	}()

	go servers.stop(servers.HttpServer)
	if servers.HttpsServer != nil {
		go servers.stop(servers.HttpsServer)
	}
}

func (servers Servers) stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), servers.ShutdownTimeout)
	defer cancel()
	server.Shutdown(ctx)
	server.Close()
}

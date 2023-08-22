package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"webserver/config"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
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
	Logger *zap.SugaredLogger
	Config *config.Config

	ShutdownTimeout time.Duration
	Errors          chan error

	HttpServer  *http.Server
	HttpsServer *http.Server

	wg *sync.WaitGroup
}

func NewServers(config *config.Config, logger *zap.SugaredLogger) Servers {
	useHTTPS := config.HTTPS()
	requestsLogger := logger.Named("requests")
	httpHandler := makeHttpHandler(config, requestsLogger, false)
	httpsHandler := makeHttpHandler(config, requestsLogger, true)

	servers := Servers{
		Logger:          logger,
		Config:          config,
		ShutdownTimeout: defaultTimeout,
		Errors:          make(chan error),
		wg:              &sync.WaitGroup{},
	}

	if useHTTPS {
		manager := autocert.Manager{
			Cache:      jsonFileCache(config.CertsFile),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Hosts()...),
		}
		httpHandler = manager.HTTPHandler(httpHandler)

		servers.HttpsServer = &http.Server{
			Addr:      fmt.Sprintf("0.0.0.0:%d", config.Ports.HTTPS),
			Handler:   httpsHandler,
			TLSConfig: manager.TLSConfig(),
		}
	}

	servers.HttpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Ports.HTTP),
		Handler: httpHandler,
	}

	return servers
}

// asynchronously starting servers
func (servers Servers) Start() {
	servers.Logger.Infoln("starting servers")

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
	servers.Logger.Infoln("stoping servers")

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

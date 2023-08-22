package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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

func NewServers(config *config.Config, logger *zap.SugaredLogger) []*http.Server {
	useHTTPS := config.HTTPS()
	servers := make([]*http.Server, 0, 2)

	logger = logger.Named("requests")
	httpHandler := makeHttpHandler(config, logger, false)
	httpsHandler := makeHttpHandler(config, logger, true)

	if useHTTPS {
		manager := autocert.Manager{
			Cache:      jsonFileCache(config.CertsFile),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.Hosts()...),
		}
		httpHandler = manager.HTTPHandler(httpHandler)

		server := &http.Server{
			Addr:      "0.0.0.0:443",
			Handler:   httpsHandler,
			TLSConfig: manager.TLSConfig(),
		}
		servers = append(servers, server)
	}

	server := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: httpHandler,
	}
	servers = append(servers, server)

	return servers
}

func StartServers(servers []*http.Server, logger *zap.SugaredLogger, shutdownTimeout time.Duration) (func(), chan error) {
	errors := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(servers))
	go func() {
		wg.Wait()
		close(errors)
	}()

	for _, server := range servers {
		go func(server *http.Server) {
			defer wg.Done()
			err := server.ListenAndServe()
			if err == http.ErrServerClosed {
				err = nil
			} else {
				err = ServerError{server, err}
			}
			errors <- err
		}(server)
	}

	stop := func() {
		logger.Infoln("stopping servers")
		for _, server := range servers {
			go func(server *http.Server) {
				ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
				defer cancel()
				server.Shutdown(ctx)
				server.Close()
			}(server)
		}
	}

	return stop, errors
}

func Wait(stop func()) {
	defer stop()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

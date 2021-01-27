package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func StartHttpsServer(certsDir string, hosts []string, port int, handler http.Handler) func() {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache(certsDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
	}
	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%d", port),
		Handler:   handler,
		TLSConfig: manager.TLSConfig(),
	}

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		server.Close()
		log.Printf("HTTPS server closed")
	}
}

func StartHttpRedirectServer(port int) func() {
	server := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := *r.URL
			url.Host = r.Host
			url.Scheme = "https"
			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		}),
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP redirect server error: %v", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		server.Close()
		log.Printf("HTTP redirect server closed")
	}
}

func StartHttpServer(port int, handler http.Handler) func() {
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		server.Close()
		log.Printf("HTTP server closed")
	}
}

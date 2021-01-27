package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// HelloWorld is a sample handler
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

// NewRouter returns a new HTTP handler that implements the main server routes
func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Set up our middleware with sane defaults
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Set up our root handlers
	router.Get("/", HelloWorld)

	// Set up our API
	//router.Mount("/api/v1/", v1.NewRouter())

	// Set up static file serving
	staticPath, _ := filepath.Abs("../../static/")
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/*", fs)

	return router
}

package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// DefaultShutdownTimeout sets the maximum amount of time to wait for the server to shutdown gracefully.
const DefaultShutdownTimeout = 10 * time.Second

// Application is contains all required base components for building web applications.
type Application struct {
	mux *chi.Mux
}

// New creates an Application that handles a set of routes for the application.
func New() *Application {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return &Application{
		mux: mux,
	}
}

// Method adds the route `pattern` that matches `method` http method to
// execute the `handler` http.Handler.
func (app *Application) Method(method, pattern string, h http.Handler) {
	app.mux.Method(method, pattern, h)
}

// ServeHTTP implements the http.Handler interface.
func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}

// Run is called to start the web service.
func (app *Application) Run() error {
	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error in ListenAndServe: %w", err)
	case <-shutdown:
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err == nil {
			return nil
		}

		// If there was an error when shutting down the server (or it timed out)
		// then we have to force it to stop.
		if err := server.Close(); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

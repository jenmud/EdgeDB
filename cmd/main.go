package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jenmud/edgedb/cmd/v1/api"
	"github.com/jenmud/edgedb/cmd/v1/web"
	_ "github.com/jenmud/edgedb/docs"
	"github.com/jenmud/edgedb/internal/server"
	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/internal/store/sqlite"
	_ "github.com/joho/godotenv/autoload"
	httpSwagger "github.com/swaggo/http-swagger"
)

// setupLogging configures the logging settings based on environment variables.
func setupLogging() {
	level := slog.LevelInfo

	switch strings.ToUpper(os.Getenv("EDGEDB_LOG_LEVEL")) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	handlerOpts := slog.HandlerOptions{
		Level:     level,
		AddSource: os.Getenv("EDGEDB_LOG_ADD_SOURCES") == "true",
	}

	var handler slog.Handler

	switch strings.ToUpper(os.Getenv("EDGEDB_LOG_HANDLER")) {
	case "JSON":
		handler = slog.NewJSONHandler(os.Stdout, &handlerOpts)
	case "TEXT":
		fallthrough
	default:
		handler = slog.NewTextHandler(os.Stdout, &handlerOpts)
	}

	// Set up structured logging with slog
	logger := slog.New(handler)

	// You can enhance this to read log level from environment variables if needed

	slog.SetDefault(logger)
}

// gracefulShutdown handles OS interrupt signals to gracefully shut down the server.
func gracefulShutdown(ctx context.Context, apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("reason", err.Error()))
	}

	slog.Info("Server exiting")
	done <- true
}

// corsMiddleware adds CORS headers to the HTTP responses.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// setupRoutes sets up all the necessary routes used by the server.
func setupRoutes(mux *http.ServeMux, s store.Store) http.Handler {
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	web.StaticAssets(mux)
	api.GETNodes(mux, s)
	api.PUTNodes(mux, s)
	api.GETEdges(mux, s)
	api.PUTEdges(mux, s)
	api.Upload(mux, s)
	return corsMiddleware(mux)
}

// @Title EdgeDB API
// @Version 1.0
// @Description EdgeDB API server
// @BasePath /
func main() {
	setupLogging()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dns := os.Getenv("EDGEDB_STORE_DSN")
	if dns == "" {
		panic("EDGEDB_STORE_DSN environment variable is not set, eg: :memory: or ./edgedb.db")
	}

	store, err := sqlite.New(ctx, dns)
	if err != nil {
		panic(fmt.Sprintf("setting up the store error: %s", err))
	}

	defer store.Close()

	mux := setupRoutes(http.NewServeMux(), store)
	server := server.NewServer(mux, os.Getenv("EDGEDB_WEB_ADDRESS"), store)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(ctx, server, done)

	slog.Info("Starting server", slog.Group("server", slog.String("address", server.Addr)))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")
}

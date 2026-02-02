package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jenmud/edgedb/internal/server"
	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/internal/store/sqlite"
	_ "github.com/joho/godotenv/autoload"
)

// setupStore sets up a new store and run the migrations.
// Defaults to a in-momory SQLite store.
func setupStore(ctx context.Context) (*store.DB, error) {
	dsn := os.Getenv("EDGEDB_STORE_DSN")
	if dsn == "" {
		dsn = ":memory:"
	}

	driver := strings.ToLower(os.Getenv("EDGEDB_STORE_DRIVER"))
	slog.SetDefault(
		slog.With(
			slog.Group(
				"store",
				slog.String("driver", driver),
				slog.String("dsn", dsn),
			),
		),
	)

	switch strings.ToLower(os.Getenv("EDGEDB_STORE_DRIVER")) {
	case "duckdb":
		return nil, errors.New("duckdb not store implemented")
	case "sqlite":
		db := sqlite.New(dsn)
		slog.Info("applying db migrations")
		return db, sqlite.ApplyMigrations(ctx, db)
	}

	return nil, errors.New("unsupported store")
}

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

// main is the entry point of the application.
func main() {
	setupLogging()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := setupStore(ctx)
	if err != nil {
		slog.Error("error setting up store", slog.String("reason", err.Error()))
		panic(fmt.Sprintf("setting up the store error: %s", err))
	}

	server := server.NewServer(os.Getenv("EDGEDB_WEB_ADDRESS"), db)

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

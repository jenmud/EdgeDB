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

	"github.com/jenmud/edgedb/internal/server"
	_ "github.com/joho/godotenv/autoload"
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
		handler = slog.NewJSONHandler(os.Stdout, &handlerOpts)
	}

	// Set up structured logging with slog
	logger := slog.New(handler)

	// You can enhance this to read log level from environment variables if needed

	slog.SetDefault(logger)
}

// gracefulShutdown handles OS interrupt signals to gracefully shut down the server.
func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	server := server.NewServer(os.Getenv("EDGEDB_WEB_ADDRESS"))

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")
}

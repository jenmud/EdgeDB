package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const defaultAddress = ":8080"

type Server struct {
	address string
}

// NewServer creates and configures a new HTTP server.
// If the address is not provided, it will default to envvar "EDGEDB_WEB_ADDRESS", or ":8080" if the envvar is not set.
func NewServer(address string) *http.Server {
	if address == "" {
		address = os.Getenv("EDGEDB_WEB_ADDRESS")
		if address == "" {
			address = defaultAddress
		}
	}

	NewServer := &Server{
		address: address,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    NewServer.address,
		Handler: NewServer.RegisterRoutes(),
		BaseContext: func(l net.Listener) context.Context {
			logger := slog.With(slog.Group("server", slog.String("address", l.Addr().String())))
			slog.SetDefault(logger)
			return context.Background()
		},
		//IdleTimeout:  time.Minute,
		//ReadTimeout:  10 * time.Second,
		//WriteTimeout: 30 * time.Second,
	}

	return server
}

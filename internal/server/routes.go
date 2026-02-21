package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jenmud/edgedb/cmd/web"
	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/models"
)

// RegisterRoutes sets up the HTTP routes for the server.
func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(web.Static))
	mux.Handle("/static/", fileServer)

	// /api/v1/nodes?term=...&limit=1000
	slog.Info("registered route", slog.String("route", "/api/v1/nodes"), slog.String("query-params", strings.Join([]string{"term", "limit"}, ",")))
	mux.HandleFunc("GET /api/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		term := strings.Trim(r.URL.Query().Get("term"), "\"")
		limit := 1000

		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
			limit = l
		}

		var (
			nodes []models.Node
			err   error
		)

		if term == "" {
			nodes, err = s.store.Nodes(ctx, store.NodesArgs{Limit: limit})
		} else {
			args := store.NodesTermSearchArgs{Term: term, Limit: limit}
			nodes, err = s.store.NodesTermSearch(ctx, args)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(nodes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

// corsMiddleware adds CORS headers to the HTTP responses.
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
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

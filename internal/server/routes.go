package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jenmud/edgedb/cmd/web"
	_ "github.com/jenmud/edgedb/docs"
	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/models"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterRoutes sets up the HTTP routes for the server.
func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fileServer := http.FileServer(http.FS(web.Static))
	mux.Handle("/static/", fileServer)

	// /api/v1/nodes?term=...&limit=1000
	slog.Info("registered route", slog.String("route", "/api/v1/nodes"), slog.String("query-params", strings.Join([]string{"term", "limit"}, ",")))
	mux.HandleFunc("PUT /api/v1/nodes", s.PUTNodes)
	mux.HandleFunc("GET /api/v1/nodes", s.GETNodes)
	mux.HandleFunc("PUT /api/v1/edges", s.PUTEdges)
	mux.HandleFunc("GET /api/v1/edges", s.GETEdges)

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

// GetNodes searches and return nodes
// @Summary Search and return nodes
// @Description Search and return nodes.
// @Tags nodes
// @Produce json
// @Param term query string false "search term" default()
// @Param snippetStart query string false "snippet start" default(<span class="text-red-500">)
// @Param snippetEnd query string false "snippet start" default(</span>)
// @Param tokens query int false "snippet tokens" minimum(1) maximum(64) default(10)
// @Param limit query int false "limit results returned" minimum(1) default(1000)
// @Success 200 {array} models.Node "List of nodes"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/nodes [get]
func (s *Server) GETNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	term := strings.Trim(r.URL.Query().Get("term"), "\"")
	snippetStart := r.URL.Query().Get("snippetStart")
	snippetEnd := r.URL.Query().Get("snippetEnd")

	limit := 1000
	tokens := 10

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}

	if s, err := strconv.Atoi(r.URL.Query().Get("tokens")); err == nil {
		tokens = s
	}

	var (
		nodes []models.Node
		err   error
	)

	if term == "" {
		nodes, err = s.store.Nodes(ctx, store.NodesArgs{Limit: limit})
	} else {
		args := store.TermSearchArgs{Term: term, Limit: limit, SnippetStart: snippetStart, SnippetEnd: snippetEnd, SnippetTokens: tokens}
		nodes, err = s.store.NodesTermSearch(ctx, args)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(nodes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUTNodesReq is a request for adding/updating one or more nodes.
type PUTNodesReq struct {
	Nodes []models.Node
}

// PUTNodes adds/update one or more nodes
// @Summary Add/update one or more nodes.
// @Description Add/update on or more nodes.
// @Tags nodes
// @Produce json
// @Param nodes body PUTNodesReq true "One or more nodes to add/update"
// @Success 200 {array} models.Node "List of nodes"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/nodes [put]
func (s *Server) PUTNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := PUTNodesReq{}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nodes, err := s.store.UpsertNodes(ctx, req.Nodes...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(nodes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetEdges searches and return edges
// @Summary Search and return edges
// @Description Search and return edges.
// @Tags edges
// @Produce json
// @Param term query string false "search term" default()
// @Param snippetStart query string false "snippet start" default(<span class="text-red-500">)
// @Param snippetEnd query string false "snippet start" default(</span>)
// @Param tokens query int false "snippet tokens" minimum(1) maximum(64) default(10)
// @Param limit query int false "limit results returned" minimum(1) default(1000)
// @Success 200 {array} models.Edge "List of edges"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/edges [get]
func (s *Server) GETEdges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	term := strings.Trim(r.URL.Query().Get("term"), "\"")
	snippetStart := r.URL.Query().Get("snippetStart")
	snippetEnd := r.URL.Query().Get("snippetEnd")

	limit := 1000
	tokens := 10

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}

	if s, err := strconv.Atoi(r.URL.Query().Get("tokens")); err == nil {
		tokens = s
	}

	var (
		edges []models.Edge
		err   error
	)

	if term == "" {
		edges, err = s.store.Edges(ctx, store.EdgesArgs{Limit: limit})
	} else {
		args := store.TermSearchArgs{Term: term, Limit: limit, SnippetStart: snippetStart, SnippetEnd: snippetEnd, SnippetTokens: tokens}
		edges, err = s.store.EdgesTermSearch(ctx, args)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(edges); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUTEdgesReq is a request for adding/updating one or more edges.
type PUTEdgesReq struct {
	Edges []models.Edge
}

// PUTEdges adds/update one or more edges
// @Summary Add/update one or more edges.
// @Description Add/update on or more edges.
// @Tags edges
// @Produce json
// @Param nodes body PUTEdgesReq true "One or more nodes to add/update"
// @Success 200 {array} models.Edge "List of edges"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/edges [put]
func (s *Server) PUTEdges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := PUTEdgesReq{}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	edges, err := s.store.UpsertEdges(ctx, req.Edges...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(edges); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

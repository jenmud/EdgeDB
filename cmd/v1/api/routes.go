package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/models"
)

// GetNodes searches and return nodes
// @Summary Search and return nodes
// @Description Search and return nodes
// @Tags nodes
// @Produce json
// @Param term query string false "search term" default()
// @Param snippetStart query string false "snippet start" default(<span class="text-red-500">)
// @Param snippetEnd query string false "snippet start" default(</span>)
// @Param tokens query int false "snippet tokens" minimum(1) maximum(64) default(10)
// @Param limit query int false "limit results returned" minimum(1) default(1000)
// @Param lastID query int false "last known ID from previous result used for pagination" default(0)
// @Success 200 {array} models.Node "List of nodes"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/nodes [get]
func GETNodes(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /api/v1/nodes"))
	mux.HandleFunc("GET /api/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		term := strings.Trim(r.URL.Query().Get("term"), "\"")
		snippetStart := r.URL.Query().Get("snippetStart")
		snippetEnd := r.URL.Query().Get("snippetEnd")

		limit := 1000
		tokens := 10
		var lastID uint64 = 0

		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
			limit = l
		}

		if s, err := strconv.Atoi(r.URL.Query().Get("tokens")); err == nil {
			tokens = s
		}

		if l, err := strconv.ParseUint(r.URL.Query().Get("lastID"), 10, 64); err == nil {
			lastID = l
		}

		var (
			nodes []models.Node
			err   error
		)

		if term == "" {
			nodes, err = s.Nodes(ctx, store.NodesArgs{Limit: limit, LastID: lastID})
		} else {
			args := store.TermSearchArgs{Term: term, Limit: limit, LastID: lastID, SnippetStart: snippetStart, SnippetEnd: snippetEnd, SnippetTokens: tokens}
			nodes, err = s.NodesTermSearch(ctx, args)
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
	})
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
func PUTNodes(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "PUT /api/v1/nodes"))
	mux.HandleFunc("PUT /api/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := PUTNodesReq{}
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nodes, err := s.UpsertNodes(ctx, req.Nodes...)
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
	})
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
// @Param lastID query int false "last known ID from previous result used for pagination" default(0)
// @Success 200 {array} models.Edge "List of edges"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/edges [get]
func GETEdges(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /api/v1/edges"))
	mux.HandleFunc("GET /api/v1/edges", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		term := strings.Trim(r.URL.Query().Get("term"), "\"")
		snippetStart := r.URL.Query().Get("snippetStart")
		snippetEnd := r.URL.Query().Get("snippetEnd")

		limit := 1000
		tokens := 10
		var lastID uint64 = 0

		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
			limit = l
		}

		if s, err := strconv.Atoi(r.URL.Query().Get("tokens")); err == nil {
			tokens = s
		}

		if l, err := strconv.ParseUint(r.URL.Query().Get("lastID"), 10, 64); err == nil {
			lastID = l
		}

		var (
			edges []models.Edge
			err   error
		)

		if term == "" {
			edges, err = s.Edges(ctx, store.EdgesArgs{Limit: limit, LastID: lastID})
		} else {
			args := store.TermSearchArgs{Term: term, Limit: limit, LastID: lastID, SnippetStart: snippetStart, SnippetEnd: snippetEnd, SnippetTokens: tokens}
			edges, err = s.EdgesTermSearch(ctx, args)
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
	})
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
func PUTEdges(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "PUT /api/v1/edges"))
	mux.HandleFunc("PUT /api/v1/edges", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := PUTEdgesReq{}
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		edges, err := s.UpsertEdges(ctx, req.Edges...)
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
	})
}

// GetGraph search and return nodes and edges used for a force directed graph
// @Summary search return nodes and edges used for a force directed graph
// @Description search return nodes and edges in a format that can be used in a force directed graph
// @Tags graph
// @Produce json
// @Param term query string false "search term" default()
// @Param snippetStart query string false "snippet start" default(<span class="text-red-500">)
// @Param snippetEnd query string false "snippet start" default(</span>)
// @Param tokens query int false "snippet tokens" minimum(1) maximum(64) default(10)
// @Param limit query int false "limit results returned" minimum(1) default(1000)
// @Success 200 {object} models.Graph "Payload used for drawing graphs."
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/graph [get]
func GETGraph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /api/v1/graph"))
	mux.HandleFunc("GET /api/v1/graph", func(w http.ResponseWriter, r *http.Request) {
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

		graph, err := s.Graph(ctx, store.TermSearchArgs{Limit: limit, Term: term, SnippetTokens: tokens, SnippetStart: snippetStart, SnippetEnd: snippetEnd})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(graph); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// PUTGraph uploads a graph using a upsert strategy.
// @Summary Uploads a graph using a upsert strategy.
// @Description Uploads a graph using a upsert strategy.
// @Tags graph
// @Produce json
// @Param nodes body models.Graph true "Graph that you are uploading"
// @Success 200 {object} models.Graph "Uploaded graph."
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/graph [put]
func PUTGraph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "PUT /api/v1/graph"))
	mux.HandleFunc("PUT /api/v1/graph", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := models.Graph{}
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nodes, err := s.UpsertNodes(ctx, req.Nodes...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		edges, err := s.UpsertEdges(ctx, req.Edges...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		resp := models.Graph{
			Nodes: nodes,
			Edges: edges,
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// GETSubGraphByNode returns a sub-graph from a node.
// @Summary Returns a graph from a node.
// @Description Returns a graph from a node.
// @Tags graph
// @Produce json
// @Param id path int true "Node id"
// @Param limit query int false "how many levels deep to return" minimum(1) default(1)
// @Success 200 {object} models.Graph "Payload used for drawing graphs."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /api/v1/graph/nodes/{id} [get]
func GETSubGraphByNode(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /api/v1/graph/nodes/{id}"))
	mux.HandleFunc("GET /api/v1/graph/nodes/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		args := store.SubGraphArgs{Limit: 1}

		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
			args.Limit = l
		}

		idstr := r.PathValue("id")
		id, err := strconv.ParseUint(idstr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		args.FromNodeID = id
		args.ToNodeID = id

		graph, err := s.SubGraph(ctx, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(graph); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// HealthStatus returns returns the health status.
// @Summary Returns returns the health status.
// @Description Returns returns the health status.
// @Tags Health
// @Produce json
// @Success 200 {object} models.Health "Current health status"
// @Failure 500 "Internal server error"
// @Router /healthz [get]
func HealthStatus(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /healthz"))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(s.Health(ctx)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

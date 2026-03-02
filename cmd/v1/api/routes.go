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
func GETNodes(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /api/v1/nodes"))
	mux.HandleFunc("GET /api/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
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
			nodes, err = s.Nodes(ctx, store.NodesArgs{Limit: limit})
		} else {
			args := store.TermSearchArgs{Term: term, Limit: limit, SnippetStart: snippetStart, SnippetEnd: snippetEnd, SnippetTokens: tokens}
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
			edges, err = s.Edges(ctx, store.EdgesArgs{Limit: limit})
		} else {
			args := store.TermSearchArgs{Term: term, Limit: limit, SnippetStart: snippetStart, SnippetEnd: snippetEnd, SnippetTokens: tokens}
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

type UploadReq struct {
	PUTNodesReq
	PUTEdgesReq
}

type UploadedResp struct {
	Nodes []models.Node
	Edges []models.Edge
}

// Upload uploads one or more nodes and edges sets.
// @Summary Uploads one or more nodes and edges sets.
// @Description Uploads one or more nodes and edges sets.
// @Tags upload
// @Produce json
// @Param nodes body UploadReq true "One or more nodes to add/update"
// @Success 200 {array} UploadedResp "List of uploaded nodes and edges"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/upload [put]
func Upload(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "PUT /api/v1/upload"))
	mux.HandleFunc("PUT /api/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := UploadReq{}
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

		resp := UploadedResp{
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

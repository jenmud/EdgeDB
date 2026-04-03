package web

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jenmud/edgedb/cmd/v1/web/view/components"
	"github.com/jenmud/edgedb/cmd/v1/web/view/pages"
	"github.com/jenmud/edgedb/internal/store"
	"github.com/starfederation/datastar-go/datastar"
)

// Static serves up static files
// @Summary Static serves up static files
// @Description Static serves up static files
// @Router /static [get]
func StaticAssets(mux *http.ServeMux) {
	slog.Info("registered route", slog.String("route", "GET /static"))

	// Access the embedded static files (using fs.Sub to get the "static" subfolder)
	sub, err := fs.Sub(Static, "static")
	if err != nil {
		panic(err)
	}

	// Create a file server to serve files from the "static" subdirectory
	fileServer := http.FileServer(http.FS(sub))

	// Handler for static assets with cache control headers
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set caching headers for static assets
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable") // Cache for 1 year
		fileServer.ServeHTTP(w, r)                                             // Serve the static file
	})))

}

func Index(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1"))
	mux.HandleFunc("GET /ui/v1", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// try and fetch the fist node in the graph, we need this because we do not
		// know what the first ID is and can not assume it starts at 1. So we
		// will fetch the first node in the table.
		nodes, _ := s.Nodes(r.Context(), store.NodesArgs{Limit: 1})
		if len(nodes) > 0 {
			http.Redirect(w, r, fmt.Sprintf("/ui/v1/graph/nodes/%d", nodes[0].ID), http.StatusMovedPermanently)
			return
		}

		pages.Index().Render(ctx, w)
	})
}

func Graph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph"))
	mux.HandleFunc("GET /ui/v1/graph", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		component := pages.GraphPage("/api/v1/graph")
		component.Render(ctx, w)
	})
}

func GraphSearch(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/search/graph"))
	mux.HandleFunc("GET /ui/v1/search/graph", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type Store struct {
			Limit int    `json:"limit"`
			Term  string `json:"term"`
		}

		queryStore := Store{}

		if err := datastar.ReadSignals(r, &queryStore); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := components.GraphContent(fmt.Sprintf("/api/v1/graph?term=%s&limit=%d", queryStore.Term, queryStore.Limit))
		component.Render(ctx, w)
	})
}

func GraphTable(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/table"))
	mux.HandleFunc("GET /ui/v1/table", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type Store struct {
			Limit int    `json:"limit"`
			Term  string `json:"term"`
		}

		queryStore := Store{}

		if err := datastar.ReadSignals(r, &queryStore); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		graph, err := s.Graph(ctx, store.TermSearchArgs{Term: queryStore.Term, Limit: queryStore.Limit})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		graphAPIUrl := fmt.Sprintf("/api/v1/graph?term=%s&limit=%d", queryStore.Term, queryStore.Limit)
		component := pages.GraphTablePage(graphAPIUrl, graph)
		component.Render(ctx, w)
	})
}

func TableSearch(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/search/table"))
	mux.HandleFunc("GET /ui/v1/search/table", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type Store struct {
			Limit int    `json:"limit"`
			Term  string `json:"term"`
		}

		queryStore := Store{}

		if err := datastar.ReadSignals(r, &queryStore); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		graph, err := s.Graph(ctx, store.TermSearchArgs{Term: queryStore.Term, Limit: queryStore.Limit})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := components.GraphTable(graph)
		component.Render(ctx, w)
	})
}

func SubGraph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph/nodes/{id}"))
	mux.HandleFunc("GET /ui/v1/graph/nodes/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		n, err := s.Node(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// component := pages.GraphPage("/api/v1/graph/nodes/" + r.PathValue("id"))
		component := pages.NodeDetailPage(n)
		component.Render(ctx, w)
	})
}

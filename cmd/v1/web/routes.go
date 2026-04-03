package web

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jenmud/edgedb/cmd/v1/web/view/pages"
	"github.com/jenmud/edgedb/internal/store"
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

// Index is the main landing page.
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

func SubGraph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph/nodes/{id}"))
	mux.HandleFunc("GET /ui/v1/graph/nodes/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// FIXME: there should be a call to query for a Node which will return a SubGraph
		n, err := s.Node(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.NodeDetailPage(n)
		component.Render(ctx, w)
	})
}

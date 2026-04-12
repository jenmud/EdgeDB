package web

import (
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jenmud/edgedb/cmd/v1/web/view/pages"
	"github.com/jenmud/edgedb/internal/store"
	"github.com/jenmud/edgedb/models"
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

// Index is the main landing page.
func Index(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1"))
	mux.HandleFunc("GET /ui/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/v1/graph/filter/table", http.StatusMovedPermanently)
		return
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

		graph, err := s.SubGraph(ctx, store.SubGraphArgs{FromNodeID: id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.NodeDetailPage(graph)
		component.Render(ctx, w)
	})
}

func FilterGraph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph/filter"))
	mux.HandleFunc("GET /ui/v1/graph/filter", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		component := pages.FilterPage()
		component.Render(ctx, w)
	})
}

func FilterGraphContent(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph/filter/content"))
	mux.HandleFunc("GET /ui/v1/graph/filter/content", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type SignalStore struct {
			Term   string
			Limit  int
			Tokens int
		}

		signals := SignalStore{}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		graph, err := s.Graph(ctx, store.TermSearchArgs{
			Term:          signals.Term,
			Limit:         signals.Limit,
			SnippetTokens: signals.Limit,
			SnippetStart:  `<span class="text-red-500">`,
			SnippetEnd:    "</span>",
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.FilterPageContent(graph)
		component.Render(ctx, w)
	})
}

func FilterGraphTable(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph/filter/table"))
	mux.HandleFunc("GET /ui/v1/graph/filter/table", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		nodes := []models.Node{}
		component := pages.FilterTablePage(nodes...)
		component.Render(ctx, w)
	})
}

func FilterGraphTableContent(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph/filter/table/content"))
	mux.HandleFunc("GET /ui/v1/graph/filter/table/content", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type SignalStore struct {
			Term   string `json:"term"`
			Limit  int    `json:"limit"`
			Count  int    `json:"count"`
			LastID uint64 `json:"lastID"`
		}

		signals := SignalStore{}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if signals.Term == "" {
			return
		}

		// FIXME: the pagenation is not working as intended, so fixing the limit for now
		//        need to also add back in the lastID when I figure out this pagenation stuff
		nodes, err := s.NodesTermSearch(ctx, store.TermSearchArgs{Term: signals.Term, Limit: 100000000000000})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.FilterTable(nodes...)
		component.Render(ctx, w)
	})
}

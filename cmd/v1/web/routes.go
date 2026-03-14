package web

import (
	"io/fs"
	"log/slog"
	"net/http"

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

	sub, err := fs.Sub(Static, "static")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(sub))
	mux.Handle("/ui/v1/static/", http.StripPrefix("/ui/v1/static/", fileServer))
}

func Index(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1"))
	mux.HandleFunc("GET /ui/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/v1/graph", http.StatusMovedPermanently)
	})
}

func Graph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph"))
	mux.HandleFunc("GET /ui/v1/graph", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		graph, err := s.Graph(ctx, store.TermSearchArgs{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.GraphPage(graph)
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

		graph, err := s.Graph(ctx, store.TermSearchArgs{Term: queryStore.Term, Limit: queryStore.Limit})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.GraphContent(graph)
		component.Render(ctx, w)
	})
}

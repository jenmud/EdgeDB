package web

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/jenmud/edgedb/cmd/v1/web/view/layout"
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
	slog.Info("registered route", slog.String("route", "GET /ui/v1/index"))
	mux.HandleFunc("GET /ui/v1/index", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		component := layout.Base("EdgeDB")
		component.Render(ctx, w)
	})
}

func Nodes(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/nodes"))
	mux.HandleFunc("GET /ui/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		nodes, err := s.Nodes(ctx, store.NodesArgs{Limit: 1000})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.NodesPage(nodes...)
		component.Render(ctx, w)
	})
}

func NodesSearch(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/search/nodes"))
	mux.HandleFunc("GET /ui/v1/search/nodes", func(w http.ResponseWriter, r *http.Request) {
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

		nodes, err := s.NodesTermSearch(ctx, store.TermSearchArgs{Limit: queryStore.Limit, Term: queryStore.Term})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.NodesTable(nodes...)
		component.Render(ctx, w)
	})
}

func Edges(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/edges"))
	mux.HandleFunc("GET /ui/v1/edges", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		edges, err := s.Edges(ctx, store.EdgesArgs{Limit: 1000})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.EdgesPage(edges...)
		component.Render(ctx, w)
	})
}

func EdgesSearch(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/search/edges"))
	mux.HandleFunc("GET /ui/v1/search/edges", func(w http.ResponseWriter, r *http.Request) {
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

		edges, err := s.EdgesTermSearch(ctx, store.TermSearchArgs{Limit: queryStore.Limit, Term: queryStore.Term})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := pages.EdgesTable(edges...)
		component.Render(ctx, w)
	})
}

func Graph(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /ui/v1/graph"))
	mux.HandleFunc("GET /ui/v1/graph", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//graph, err := s.Graph(ctx, store.TermSearchArgs{})
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		component := pages.GraphPage()
		component.Render(ctx, w)
	})
}

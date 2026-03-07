package web

import (
	"log/slog"
	"net/http"

	"github.com/jenmud/edgedb/cmd/v1/web/view/layout"
	"github.com/jenmud/edgedb/internal/store"
)

// Static serves up static files
// @Summary Static serves up static files
// @Description Static serves up static files
// @Router /static [get]
func StaticAssets(mux *http.ServeMux) {
	slog.Info("registered route", slog.String("route", "GET /static"))
	fileServer := http.FileServer(http.FS(Static))
	mux.Handle("/static/", fileServer)
}

func Index(mux *http.ServeMux, s store.Store) {
	slog.Info("registered route", slog.String("route", "GET /v1/ui"))
	mux.HandleFunc("GET /v1/ui", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		component := layout.Base("EdgeDB")
		component.Render(ctx, w)
	})
}

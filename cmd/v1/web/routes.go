package web

import (
	"log/slog"
	"net/http"
)

// Static serves up static files
// @Summary Static serves up static files
// @Description Static serves up static files
// @Router /api/v1/static [get]
func StaticAssets(mux *http.ServeMux) {
	slog.Info("registered route", slog.String("route", "PUT /api/v1/static"))
	fileServer := http.FileServer(http.FS(Static))
	mux.Handle("/static/", fileServer)
}

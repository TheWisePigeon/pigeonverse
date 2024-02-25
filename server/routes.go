package server

import (
	"embed"
	"net/http"
	"pigeonverse/handlers"
)

func RegisterRoutes(mux *http.ServeMux, fs *embed.FS) {
	mux.Handle("GET /", handlers.RenderLandingPage(fs))
}

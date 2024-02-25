package server

import (
	"embed"
	"net/http"
	"pigeonverse/handlers"
)

func RegisterRoutes(mux *http.ServeMux, fs *embed.FS) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))
	mux.Handle("GET /", handlers.RenderLandingPage(fs))
}

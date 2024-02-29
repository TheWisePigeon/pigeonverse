package server

import (
	"net/http"
	"pigeonverse/handlers"
)

func RegisterRoutes(mux *http.ServeMux, contentDir string) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))
	mux.Handle("GET /", handlers.RenderLandingPage())
	mux.Handle("GET /posts", handlers.RenderPostsPage(contentDir))
	mux.Handle("GET /posts/{slug}", handlers.RenderPost(contentDir))
}

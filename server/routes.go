package server

import (
	"embed"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, fs *embed.FS) {
	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Well, hello there"))
		return
	}))
}

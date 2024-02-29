package server

import (
	"net/http"
	"os"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	contentDir := os.Getenv("CONTENT_DIR")
	RegisterRoutes(mux, contentDir)
	return mux
}

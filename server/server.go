package server

import (
	"net/http"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}

package main

import (
	"embed"
	"log"
	"net"
	"net/http"
	"pigeonverse/server"
)

//go:embed views
var viewsFS embed.FS

func main() {
	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: server.NewServer(&viewsFS),
	}
	log.Println("Server launched")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

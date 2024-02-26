package main

import (
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"pigeonverse/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", port),
		Handler: server.NewServer(),
	}
	log.Println("Server launched")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

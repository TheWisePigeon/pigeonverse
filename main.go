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
	if os.Getenv("ENV") != "PROD" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", port),
		Handler: server.NewServer(),
	}
	log.Println("Server launched")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

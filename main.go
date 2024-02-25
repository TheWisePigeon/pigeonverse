package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed views
var viewsFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable not set")
	}
	fmt.Println("Hello world")
}

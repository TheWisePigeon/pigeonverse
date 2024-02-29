package handlers

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PostData struct {
	Title    string `json:"title"`
	PostedAt string `json:"posted_at"`
	Slug     string `json:"slug"`
	TLDR     string `json:"tldr"`
}

func getPostsData() ([]PostData, error) {
	contentDir := os.Getenv("CONTENT_DIR")
	entries, err := os.ReadDir(contentDir)
	data := []PostData{}
	if err != nil {
		return data, fmt.Errorf("Error while reading content dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := os.Open(filepath.Join(contentDir, entry.Name()))
		if err != nil {
			return data, fmt.Errorf("Error while opening mardkdown file: %w", err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		newPostInfo := PostData{}
		for i := 0; i < 6 && scanner.Scan(); i++ {
			parts := strings.SplitN(scanner.Text(), " ", 2)
			switch i {
			case 1:
				newPostInfo.Title = parts[1]
			case 2:
				newPostInfo.PostedAt = parts[1]
			case 3:
				newPostInfo.Slug = parts[1]
			case 4:
				newPostInfo.PostedAt = parts[1]
			}
		}
		data = append(data, newPostInfo)
	}
	return data, nil
}

func RenderLandingPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles("views/base.html", "views/index.html")
		if err != nil {
			log.Println("Error while parsing template", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = templ.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println("Error while rendering landing page", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func RenderPostsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := getPostsData()
		if err != nil {
			log.Println("Error while reading frontmatterd data", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		templ, err := template.ParseFiles("views/base.html", "views/posts.html")
		if err != nil {
			log.Println("Error while parsing template", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		templData := map[string]interface{}{
			"Posts": data,
		}
		err = templ.ExecuteTemplate(w, "base", templData)
		if err != nil {
			log.Println("Error while rendering posts page", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

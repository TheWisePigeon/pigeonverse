package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type PostInfo struct {
	Title    string `json:"title"`
	PostedAt string `json:"posted_at"`
	Slug     string `json:"slug"`
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
		postsData := []PostInfo{}
		entry, err := os.ReadDir(os.Getenv("CONTENT_DIR"))
		if err != nil {
			log.Println("Error while reading content dir:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, file := range entry {
			slug := strings.TrimSuffix(file.Name(), ".md")
			parts := strings.Split(slug, "-")
			title := strings.ReplaceAll(parts[0], "_", " ")
			postedAt := strings.ReplaceAll(parts[1], "_", " ")
			postData := PostInfo{
				Title:    title,
				PostedAt: postedAt,
				Slug:     slug,
			}
			postsData = append(postsData, postData)
		}
		templ, err := template.ParseFiles("views/base.html", "views/posts.html")
		if err != nil {
			log.Println("Error while parsing template", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		templData := map[string]interface{}{
			"Posts": postsData,
		}
		err = templ.ExecuteTemplate(w, "base", templData)
		if err != nil {
			log.Println("Error while rendering posts page", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

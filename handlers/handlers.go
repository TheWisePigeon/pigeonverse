package handlers

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pigeonverse/helpers"
)

type PostData struct {
	helpers.Frontmatter
	Content template.HTML
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

func RenderPostsPage(contentDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir(contentDir)
		if err != nil {
			log.Println("Error while reading content dir: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := []helpers.Frontmatter{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			filePath := filepath.Join(contentDir, entry.Name())
			postFrontmatter, _, err := helpers.ExtractFrontmatter(filePath)
			if err != nil {
				log.Println("Error while reading post frontmatter: ", err)
				return
			}
			data = append(data, *postFrontmatter)
		}
		if err != nil {
			log.Println("Error while reading frontmatter data", err)
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

func RenderPost(contentDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postData := PostData{}
		slug := r.PathValue("slug")
		postFilePath := filepath.Join(contentDir, fmt.Sprintf("%s.md", slug))
		_, err := os.Stat(postFilePath)
		if os.IsNotExist(err) {
			fmt.Println("came here")
			templ, err := template.ParseFiles("views/base.html", "views/post_not_found.html")
			if err != nil {
				log.Println("Error while parsing templates: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			templ.ExecuteTemplate(w, "base", nil)
			return
		}
		frontmatter, postBody, err := helpers.ExtractFrontmatter(postFilePath)
		if err != nil {
			log.Println("Error while extracting frontmatter: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		postData.Frontmatter = *frontmatter
		htmlContent := blackfriday.MarkdownCommon([]byte(postBody))
		templ, err := template.ParseFiles("views/base.html", "views/post.html")
		if err != nil {
			log.Println("Error while parsing template: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		postData.Content = template.HTML(htmlContent)
		err = templ.ExecuteTemplate(w, "base", postData)
		if err != nil {
			log.Println("Error while executing template: ", err.Error())
			return
		}
	})
}

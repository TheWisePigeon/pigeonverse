package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pigeonverse/helpers"
	"github.com/russross/blackfriday"
)

type Project struct {
	Name        string

	Description string
	TechStack   []string
	Github      string
	LiveDemo    string
	Finished    bool
}

var (
	projects = []Project{
		{
			Name:        "Yoot",
			Description: "Minimalist CMS designed for simplicity and ease of use",
			TechStack:   []string{"Sveltekit", "TypeScript", "Postgres"},
			Github:      "https://github.com/TheWisePigeon/yoot",
			LiveDemo:    "",
			Finished:    true,
		},
		{
			Name:        "Certus",
			Description: "HTTP test runner",
			TechStack:   []string{"Rust"},
			Github:      "https://github.com/TheWisePigeon/certus",
			LiveDemo:    "",
			Finished:    true,
		},
		{
			Name:        "Visio",
			Description: "Cloud based service that provides face detection and recognition ",
			TechStack:   []string{"Go", "Postgres", "Redis", "Docker"},
			Github:      "https://github.com/TheWisePigeon/visio",
			LiveDemo:    "https://visio-beta.onrender.com",
			Finished:    false,
		},
		{
			Name:        "Restdis",
			Description: "Postgrest equivalent for Redis",
			TechStack:   []string{"Go", "SQLite", "Docker"},
			Github:      "https://github.com/TheWisePigeon/restdis",
			LiveDemo:    "",
			Finished:    false,
		},
		{
			Name:        "SQL to TypeScript",
			Description: "Convert your SQL tables into typescript types",
			TechStack:   []string{"Go"},
			Github:      "https://github.com/TheWisePigeon/sql-to-typescript",
			LiveDemo:    "",
			Finished:    true,
		},
		{
			Name:        "Rex",
			Description: "ExpressJS project scaffolder",
			TechStack:   []string{"Rust"},
			Github:      "https://github.com/TheWisePigeon/rex",
			LiveDemo:    "",
			Finished:    true,
		},
		{
			Name:        "Tabula",
			Description: "Your place to organize your tasks with ease ",
			TechStack:   []string{"Sveltekit", "TypeScript"},
			Github:      "https://github.com/TheWisePigeon/tabula",
			LiveDemo:    "https://tabula.lol",
			Finished:    true,
		},
	}
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
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			data = append(data, *postFrontmatter)
		}
		if err != nil {
			log.Println("Error while reading frontmatter data", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		templ, err := template.ParseFiles("views/base.html", "views/posts.html")
		if err != nil {
			log.Println("Error while parsing template", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		templData := map[string]interface{}{
			"Posts": data,
		}
		err = templ.ExecuteTemplate(w, "base", templData)
		if err != nil {
			log.Println("Error while rendering posts page", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
			templ, err := template.ParseFiles("views/base.html", "views/post_not_found.html")
			if err != nil {
				log.Println("Error while parsing templates: ", err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			templ.ExecuteTemplate(w, "base", nil)
			return
		}
		frontmatter, postBody, err := helpers.ExtractFrontmatter(postFilePath)
		if err != nil {
			log.Println("Error while extracting frontmatter: ", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		postData.Frontmatter = *frontmatter
		htmlContent := blackfriday.MarkdownCommon([]byte(postBody))
		templ, err := template.ParseFiles("views/post.html")
		if err != nil {
			log.Println("Error while parsing template: ", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		postData.Content = template.HTML(htmlContent)
		err = templ.ExecuteTemplate(w, "main", postData)
		if err != nil {
			log.Println("Error while executing template: ", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	})
}

func RenderProjectsPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles("views/base.html", "views/projects.html")
		if err != nil {
			log.Println("Error while parsiong templates:", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		projectsData := map[string][]Project{
			"Projects": projects,
		}
		err = templ.ExecuteTemplate(w, "base", projectsData)
		if err != nil {
			log.Println("Error while rendering projects page:", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	})
}

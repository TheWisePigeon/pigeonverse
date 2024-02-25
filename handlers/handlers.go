package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

func RenderLandingPage(fs *embed.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFS(fs, "views/index.html")
		if err != nil {
			log.Println("Error while parsing template", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = templ.ExecuteTemplate(w, "home", nil)
		if err != nil {
			log.Println("Error while rendering page", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

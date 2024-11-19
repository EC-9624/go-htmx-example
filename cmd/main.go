package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)


func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	tmplPath := filepath.Join("htmx", "templates", tmpl)
	t, err := template.ParseFiles(
		filepath.Join("htmx", "templates", "layout.html"),
		tmplPath,
	)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the request is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Render only the content block
		if err := t.ExecuteTemplate(w, "content", data); err != nil {
			http.Error(w, "Error executing content template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Render the full layout
		if err := t.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "Error executing layout template: "+err.Error(), http.StatusInternalServerError)
		}
	}
	
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "1-tabs-navigation.html", nil)
}

func multiSelectPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "2-multi-select.html", nil)
}

func externalApiPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "3-external-api.html", nil)
}

func webSocketPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "4-web-socket.html", nil)
}

func main() {
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/multi-select", multiSelectPageHandler)
	http.HandleFunc("/external-api", externalApiPageHandler)
	http.HandleFunc("/web-socket", webSocketPageHandler)


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

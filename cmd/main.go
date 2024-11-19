package main

import (
	"log"
	"net/http"
	"path/filepath"

	"go-htmx-example/internal"
)

func main() {
	// Initialize template renderer
	renderer := internal.NewTemplateRenderer(
		filepath.Join("htmx", "templates"),
		"layout.html",
	)

	// Initialize handlers
	h := internal.NewHandlers(renderer)

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HomePage)
	
	mux.HandleFunc("/multi-select", h.MultiSelectHandler)            // Handle main page
	mux.HandleFunc("/multi-select/table/", h.HandleMultiSelectToggle)     // Handle table updates

	mux.HandleFunc("/external-api", h.ExternalApi)

	mux.HandleFunc("/web-socket", h.WebSocket)

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

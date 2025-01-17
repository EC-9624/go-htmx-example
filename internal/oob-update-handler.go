package internal

import (
	"fmt"
	"net/http"
)

// Count struct to hold the count variable
type Count struct {
    Count int
}

// Global variable count
var count = Count{Count: 0}

func (h *Handlers) OobUpdate(w http.ResponseWriter, r *http.Request) {
	
	h.renderer.Render(w, r, "5-oob-update.html", count)
}

func (h *Handlers) AddCount(w http.ResponseWriter, r *http.Request) {

    count.Count++
	h.renderer.Render(w, r, "oob-response.html", count)
}

func (h *Handlers) RemoveCount(w http.ResponseWriter, r *http.Request) {

    count.Count--
	h.renderer.Render(w, r, "oob-response.html", count)
}

// GET count on load
func (h *Handlers) GetCount(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, `<span id="oob-span" hx-get="/get-count" hx-trigger="load" hx-target="#oob-span" hx-swap="outerHTML"
        class="inline-flex items-center px-2 text-white bg-red-500 rounded-full shadow-md">
        %d
    </span>`, count.Count)
}

package internal

import (
	"net/http"
	"strings"
)

// Person represents an individual's data in the multi-select interface
type Person struct {
	Name             string
	Email            string
	Role             string
	ImageUrl         string
	LastSeen         *string
	LastSeenDateTime *string
}

// MultiSelectPageData contains the configuration and data for the multi-select page
type MultiSelectPageData struct {
	NameEnabled       bool
	EmailEnabled      bool
	LastOnlineEnabled bool
	People           []Person
}

// default data for the page
var defaultPeople = []Person{
	{
		Name:             "Leslie Alexander",
		Email:            "leslie.alexander@example.com",
		Role:            "Co-Founder / CEO",
		ImageUrl:        "https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format",
		LastSeen:        stringPtr("3h ago"),
		LastSeenDateTime: stringPtr("2023-01-23T13:23Z"),
	},
	{
		Name:             "Michael Foster",
		Email:            "michael.foster@example.com",
		Role:            "Co-Founder / CTO",
		ImageUrl:        "https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?auto=format",
		LastSeen:        stringPtr("3h ago"),
		LastSeenDateTime: stringPtr("2023-01-23T13:23Z"),
	},
	{
		Name:             "Dries Vincent",
		Email:            "dries.vincent@example.com",
		Role:            "Business Relations",
		ImageUrl:        "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?auto=format",
		LastSeen:        nil,
		LastSeenDateTime: nil,
	},
}

var currentState = MultiSelectPageData{
	NameEnabled:       true,
	EmailEnabled:      true,
	LastOnlineEnabled: true,
	People:           defaultPeople,
}

// Handle GET request - show the initial page
func (h *Handlers) MultiSelectHandler(w http.ResponseWriter, r *http.Request) {
	
	
	h.renderer.Render(w, r, "2-multi-select.html", currentState)
}

// handleMultiSelectToggle processes the POST requests for toggling fields
func (h *Handlers) HandleMultiSelectToggle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	field := parts[3]    // name, email, or last-online
	
	 // Toggle the appropriate field
	 switch field {
    case "name":
        currentState.NameEnabled = !currentState.NameEnabled
    case "email":
        currentState.EmailEnabled = !currentState.EmailEnabled
    case "last-online":
        currentState.LastOnlineEnabled = !currentState.LastOnlineEnabled
    default:
        http.Error(w, "Invalid field", http.StatusBadRequest)
        return
    }

	// Render the updated content
	h.renderer.Render(w, r, "2-multi-select.html", currentState)
}

// stringPtr is a helper function to get a pointer to a string
func stringPtr(s string) *string {
	return &s
}

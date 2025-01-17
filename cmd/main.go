package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"go-htmx-example/internal"
	"go-htmx-example/internal/hardware"
)

func main() {
	// Initialize template renderer
	renderer := internal.NewTemplateRenderer(
		filepath.Join("htmx", "templates"),
		"layout.html",
	)

	// Initialize handlers
	h := internal.NewHandlers(renderer)

	// Initialize WebSocket server
	wsServer := internal.NewWebSocketServer()

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HomePage)
	
	mux.HandleFunc("/multi-select", h.MultiSelectHandler)            // Handle main page
	mux.HandleFunc("/multi-select/table/", h.HandleMultiSelectToggle)     // Handle table updates

	mux.HandleFunc("/external-api", h.ExternalApiHandler)
	mux.HandleFunc("/poke", h.HandlePokeRequest)
	
	mux.HandleFunc("/web-socket", h.WebSocket) // Handle websocket page
	mux.HandleFunc("/oob-update", h.OobUpdate)
	mux.HandleFunc("/count", h.AddCount)
	mux.HandleFunc("/remove-count", h.RemoveCount)
	mux.HandleFunc("/get-count", h.GetCount)





	// WebSocket route
	mux.HandleFunc("/ws", wsServer.SubscribeHandler)
	mux.HandleFunc("/ws-json", wsServer.SubscribeHandlerJson) 

	// Start hardware monitoring goroutine
	go startMonitoring(wsServer)

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

func startMonitoring(wsServer *internal.WebSocketServer) {
	for {
		// Fetch hardware data
		systemData, err := hardware.GetSystemSection()
		if err != nil {
			fmt.Println("Error fetching system data:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		diskData, err := hardware.GetDiskSection()
		if err != nil {
			fmt.Println("Error fetching disk data:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		cpuData, err := hardware.GetCpuSection()
		if err != nil {
			fmt.Println("Error fetching CPU data:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Create HTML message for `/ws`
		timeStamp := time.Now().Format("2006-01-02 15:04:05")
		htmlMsg := []byte(fmt.Sprintf(`
      <div hx-swap-oob="innerHTML:#update-timestamp">
        <p><i style="color: green" class="fa fa-circle"></i> %s</p>
      </div>
      <div hx-swap-oob="innerHTML:#system-data">%s</div>
      <div hx-swap-oob="innerHTML:#cpu-data">%s</div>
      <div hx-swap-oob="innerHTML:#disk-data">%s</div>`,
			timeStamp, systemData, cpuData, diskData,
		))

		// Fetch system info JSON
		systemInfoJSON, err := hardware.GetSystemInfoJSON()
		if err != nil {
			fmt.Println("Error fetching system info:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Broadcast HTML message to `/ws` route
		wsServer.PublishMessage(htmlMsg)  
		// Broadcast JSON message to `/ws-json` route
		wsServer.PublishMessage([]byte(systemInfoJSON))  

		// Sleep for 1 second
		time.Sleep(1 * time.Second)
	}
}

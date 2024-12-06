package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WebSocketServer struct {
	SubscriberMessageBuffer int
	SubscribersMu           sync.Mutex
	Subscribers             map[*Subscriber]struct{}
}

type Subscriber struct {
	Msgs chan []byte
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		SubscriberMessageBuffer: 10,
		Subscribers:             make(map[*Subscriber]struct{}),
	}
}

func (s *WebSocketServer) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println("Subscription error:", err)
	}
}

func (s *WebSocketServer) SubscribeHandlerJson(w http.ResponseWriter, r *http.Request) {
	// Add these headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println("Subscription error:", err)
	}
}

func (s *WebSocketServer) addSubscriber(subscriber *Subscriber) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()
	s.Subscribers[subscriber] = struct{}{}
	fmt.Println("Added subscriber:", subscriber)
}

func (s *WebSocketServer) removeSubscriber(subscriber *Subscriber) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()
	delete(s.Subscribers, subscriber)
	fmt.Println("Removed subscriber:", subscriber)
}

func (s *WebSocketServer) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Explicitly set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Configure WebSocket accept options to allow specific origins
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,  // This will bypass strict origin checking
		OriginPatterns: []string{
			"http://localhost:5173",
			"https://localhost:5173",
			"localhost:5173",
			"localhost:8080",
		},
	})
	if err != nil {
		return fmt.Errorf("websocket accept error: %w", err)
	}
	defer c.Close(websocket.StatusInternalError, "Internal error")

	subscriber := &Subscriber{
		Msgs: make(chan []byte, s.SubscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)
	defer func() {
		s.removeSubscriber(subscriber)
		close(subscriber.Msgs)
	}()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-subscriber.Msgs:
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.Write(writeCtx, websocket.MessageText, msg)
			cancel() 
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *WebSocketServer) PublishMessage(msg []byte) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()

	for subscriber := range s.Subscribers {
		select {
		case subscriber.Msgs <- msg:
		default:
			fmt.Println("Subscriber channel full, skipping...")
		}
	}
}

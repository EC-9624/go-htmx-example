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
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
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
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
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

package main

import (
	"fmt"
	"log"
	"net/http"

	eventsub "github.com/dnsge/twitch-eventsub-framework/v2"
	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
)

const (
	// Key for verifying webhook requests
	secretKey = `hey this is really secret`
)

func main() {
	handler := eventsub.NewSubHandler(true, []byte(secretKey))
	handler.HandleChannelUpdate = func(h *bindings.NotificationHeaders, event *bindings.EventChannelUpdate) {
		fmt.Println("Got a channel.update notification!")
		fmt.Printf("Message id: %s\n", h.MessageID)
		fmt.Printf("Channel: %s Title: %s\n", event.BroadcasterUserName, event.Title)
	}

	// Listen on port 8080. In a real application, you would pass your mux with
	// a route that uses the handler.
	err := http.ListenAndServe("127.0.0.1:8080", handler)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
}

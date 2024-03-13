package main

import (
	"context"
	"fmt"

	"github.com/dnsge/twitch-eventsub-framework/v2"
	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
)

const (
	// These are usually created by your application automatically
	clientID = `abc123`
	appToken = `def456`
	// Key for verifying webhook requests
	secretKey = `hey this is really secret`
)

func main() {
	client := eventsub.NewSubClient(eventsub.NewStaticCredentials(clientID, appToken))
	res, _ := client.Subscribe(context.Background(), &eventsub.SubRequest{
		Type: "channel.update",
		Condition: bindings.ConditionChannelUpdate{
			BroadcasterUserID: "22484632",
		},
		Callback: "https://my.website/api/twitch/webhooks",
		Secret:   secretKey,
	})

	fmt.Printf("Using %d/%d of webhook cost limit\n", res.TotalCost, res.MaxTotalCost)
}

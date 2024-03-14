<div align="center">

<h1>twitch-eventsub-framework</h1>

Framework for Twitch EventSub applications built with webhooks

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci]

</div>

## Installation

`go get -u github.com/dnsge/twitch-eventsub-framework/v2`

Go 1.21+ is required.

## Quick Start

Use a `SubHandler` to listen for incoming notifications from Twitch servers.

```go
// Create my handler with verification and a secret key
handler := eventsub.NewSubHandler(true, []byte(`my signing secret`))

// Process channel.update EventSub notifications
handler.HandleChannelUpdate = func(h *bindings.NotificationHeaders, event *bindings.EventChannelUpdate) {
    fmt.Println("Got a channel.update notification!")
    fmt.Printf("Channel = %s, Title = %s\n", event.BroadcasterUserName, event.Title)
}

// Listen for HTTP requests from Twitch EventSub servers
http.ListenAndServe("127.0.0.1:8080", handler)
```

Use a `SubClient` to subscribe to EventSub subscriptions.

```go
// Create a client with a ClientID and App Token
client := eventsub.NewSubClient(eventsub.NewStaticCredentials(clientID, appToken))

// Subscribe to channel.update events for forsen
client.Subscribe(context.Background(), &eventsub.SubRequest{
    Type: "channel.update",
    Condition: bindings.ConditionChannelUpdate{
        BroadcasterUserID: "22484632",
    },
    Callback: "https://my.website/api/twitch/webhooks",
    Secret:   `my signing secret`,
})
```

## Examples

1. See [examples/sub_client/main.go](examples/sub_client/main.go) for an example usage of creating a new webhook subscription.
2. See [examples/sub_handler/main.go](examples/sub_handler/main.go) for an example usage of receiving webhook notifications from Twitch.

[doc-img]: https://pkg.go.dev/badge/github.com/dnsge/twitch-eventsub-framework/v2
[doc]: https://pkg.go.dev/github.com/dnsge/twitch-eventsub-framework/v2
[ci-img]: https://github.com/dnsge/twitch-eventsub-framework/actions/workflows/go.yml/badge.svg?branch=v2
[ci]: https://github.com/dnsge/twitch-eventsub-framework/actions/workflows/go.yml?branch=v2

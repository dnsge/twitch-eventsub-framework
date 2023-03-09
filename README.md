# twitch-eventsub-framework

A small framework for creating Twitch EventSub applications with an HTTP transport.

Note: At the moment, this library does not support the WebSocket transport.

## Features

This package has two main features:
1. A `SubClient` to subscribe to, unsubscribe to, and list subscriptions created with EventSub
2. A `SubHandler` to handle webhook verification requests and dispatch webhook notifications

## Examples
1. See [examples/sub_client/main.go](examples/sub_client/main.go) for an example usage of creating a new webhook subscription.
2. See [examples/sub_handler/main.go](examples/sub_handler/main.go) for an example usage of receiving webhook notifications from Twitch.

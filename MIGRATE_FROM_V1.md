# Migrating from v1

Multiple breaking changes were made to the API while creating version 2 of `twitch-eventsub-framework`.

## Renamed:
- `eventsub_framework` package renamed to `eventsub`
- `SubClient` renamed to `Client`
- `SubHandler` renamed to `Handler`

## Changed:
- Go 1.21+ required
- `NewHandler` only takes slice for secret
- `github.com/dnsge/twitch-eventsub-bindings` package replaced with local `bindings` package
- `Handler.VerifyChallenge` takes `context.Context`
- `Handler.OnDuplicateNotification` takes `context.Context`
- `Credientials` interface takes `context.Context`
- Event handlers now have the following signature: 
```go
type EventHandler[EventMessage any] func(headers bindings.NotificationHeaders, sub bindings.Subscription, event EventMessage)
```

## New:
- `TrackerFunc` wrapper helper
- `Handler` handler functions are now generated from the struct definition
- `Handler` now has `BeforeHandleEvent` to trigger logic before dispatching the appropriate event handler

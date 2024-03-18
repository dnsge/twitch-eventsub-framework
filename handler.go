package eventsub

//go:generate go run ./scripts/handler_generator --input=$GOFILE --output=notification_handler.go

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mozillazg/go-httpheader"

	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
	"github.com/dnsge/twitch-eventsub-framework/v2/bindings/beta"
)

const (
	webhookCallbackVerification = "webhook_callback_verification"
	notificationMessageType     = "notification"
)

// EventHandler is an event callback to process a notification from EventSub.
type EventHandler[EventMessage any] func(*bindings.NotificationHeaders, *EventMessage)

// Handler implements http.Handler to receive Twitch EventSub webhook requests.
//
// Handler handles both verification of new subscriptions and dispatching of
// event notifications. To handle a specific event, set the corresponding
// HandleXXX struct field. When a notification is received and validated, the
// handler function will be invoked in a new goroutine.
type Handler struct {
	// Secret used to compute signature, or nil if not enabled.
	signatureSecret []byte

	// VerifyChallenge is called to determine whether a subscription challenge
	// should be accepted.
	VerifyChallenge func(context.Context, *bindings.NotificationHeaders, *bindings.SubscriptionChallenge) bool

	// IDTracker used to deduplicate notifications
	IDTracker IDTracker
	// OnDuplicateNotification is called when the provided IDTracker rejects a
	// EventSub notification as duplicate.
	OnDuplicateNotification func(context.Context, *bindings.NotificationHeaders)

	HandleChannelUpdate                       EventHandler[bindings.EventChannelUpdate]                       `eventsub-type:"channel.update" eventsub-version:"2"`
	HandleChannelFollow                       EventHandler[bindings.EventChannelFollow]                       `eventsub-type:"channel.follow" eventsub-version:"2"`
	HandleChannelSubscribe                    EventHandler[bindings.EventChannelSubscribe]                    `eventsub-type:"channel.subscribe" eventsub-version:"1"`
	HandleChannelSubscriptionEnd              EventHandler[bindings.EventChannelSubscriptionEnd]              `eventsub-type:"channel.subscription.end" eventsub-version:"1"`
	HandleChannelSubscriptionGift             EventHandler[bindings.EventChannelSubscriptionGift]             `eventsub-type:"channel.subscription.gift" eventsub-version:"1"`
	HandleChannelSubscriptionMessage          EventHandler[bindings.EventChannelSubscriptionMessage]          `eventsub-type:"channel.subscription.message" eventsub-version:"1"`
	HandleChannelCheer                        EventHandler[bindings.EventChannelCheer]                        `eventsub-type:"channel.cheer" eventsub-version:"1"`
	HandleChannelRaid                         EventHandler[bindings.EventChannelRaid]                         `eventsub-type:"channel.raid" eventsub-version:"1"`
	HandleChannelBan                          EventHandler[bindings.EventChannelBan]                          `eventsub-type:"channel.ban" eventsub-version:"1"`
	HandleChannelUnban                        EventHandler[bindings.EventChannelUnban]                        `eventsub-type:"channel.unban" eventsub-version:"1"`
	HandleChannelModeratorAdd                 EventHandler[bindings.EventChannelModeratorAdd]                 `eventsub-type:"channel.moderator.add" eventsub-version:"1"`
	HandleChannelModeratorRemove              EventHandler[bindings.EventChannelModeratorRemove]              `eventsub-type:"channel.moderator.remove" eventsub-version:"1"`
	HandleChannelPointsRewardAdd              EventHandler[bindings.EventChannelPointsRewardAdd]              `eventsub-type:"channel.channel_points_custom_reward.add" eventsub-version:"1"`
	HandleChannelPointsRewardUpdate           EventHandler[bindings.EventChannelPointsRewardUpdate]           `eventsub-type:"channel.channel_points_custom_reward.update" eventsub-version:"1"`
	HandleChannelPointsRewardRemove           EventHandler[bindings.EventChannelPointsRewardRemove]           `eventsub-type:"channel.channel_points_custom_reward.remove" eventsub-version:"1"`
	HandleChannelPointsRewardRedemptionAdd    EventHandler[bindings.EventChannelPointsRewardRedemptionAdd]    `eventsub-type:"channel.channel_points_custom_reward_redemption.add" eventsub-version:"1"`
	HandleChannelPointsRewardRedemptionUpdate EventHandler[bindings.EventChannelPointsRewardRedemptionUpdate] `eventsub-type:"channel.channel_points_custom_reward_redemption.update" eventsub-version:"1"`
	HandleChannelPollBegin                    EventHandler[bindings.EventChannelPollBegin]                    `eventsub-type:"channel.poll.begin" eventsub-version:"1"`
	HandleChannelPollProgress                 EventHandler[bindings.EventChannelPollProgress]                 `eventsub-type:"channel.poll.progress" eventsub-version:"1"`
	HandleChannelPollEnd                      EventHandler[bindings.EventChannelPollEnd]                      `eventsub-type:"channel.poll.end" eventsub-version:"1"`
	HandleChannelPredictionBegin              EventHandler[bindings.EventChannelPredictionBegin]              `eventsub-type:"channel.prediction.begin" eventsub-version:"1"`
	HandleChannelPredictionProgress           EventHandler[bindings.EventChannelPredictionProgress]           `eventsub-type:"channel.prediction.progress" eventsub-version:"1"`
	HandleChannelPredictionLock               EventHandler[bindings.EventChannelPredictionLock]               `eventsub-type:"channel.prediction.lock" eventsub-version:"1"`
	HandleChannelPredictionEnd                EventHandler[bindings.EventChannelPredictionEnd]                `eventsub-type:"channel.prediction.end" eventsub-version:"1"`
	HandleDropEntitlementGrant                EventHandler[bindings.EventDropEntitlementGrant]                `eventsub-type:"drop.entitlement.grant" eventsub-version:"1"`
	HandleExtensionBitsTransactionCreate      EventHandler[bindings.EventBitsTransactionCreate]               `eventsub-type:"extension.bits_transaction.create" eventsub-version:"1"`
	HandleGoalBegin                           EventHandler[bindings.EventGoals]                               `eventsub-type:"channel.goal.begin" eventsub-version:"1"`
	HandleGoalProgress                        EventHandler[bindings.EventGoals]                               `eventsub-type:"channel.goal.progress" eventsub-version:"1"`
	HandleGoalEnd                             EventHandler[bindings.EventGoals]                               `eventsub-type:"channel.goal.end" eventsub-version:"1"`
	HandleHypeTrainBegin                      EventHandler[bindings.EventHypeTrainBegin]                      `eventsub-type:"channel.hype_train.begin" eventsub-version:"1"`
	HandleHypeTrainProgress                   EventHandler[bindings.EventHypeTrainProgress]                   `eventsub-type:"channel.hype_train.progress" eventsub-version:"1"`
	HandleHypeTrainEnd                        EventHandler[bindings.EventHypeTrainEnd]                        `eventsub-type:"channel.hype_train.end" eventsub-version:"1"`
	HandleStreamOnline                        EventHandler[bindings.EventStreamOnline]                        `eventsub-type:"stream.online" eventsub-version:"1"`
	HandleStreamOffline                       EventHandler[bindings.EventStreamOffline]                       `eventsub-type:"stream.offline" eventsub-version:"1"`
	HandleUserUpdate                          EventHandler[bindings.EventUserUpdate]                          `eventsub-type:"user.update" eventsub-version:"1"`
	HandleUserAuthorizationGrant              EventHandler[bindings.EventUserAuthorizationGrant]              `eventsub-type:"user.authorization.grant" eventsub-version:"1"`
	HandleUserAuthorizationRevoke             EventHandler[bindings.EventUserAuthorizationRevoke]             `eventsub-type:"user.authorization.revoke" eventsub-version:"1"`
	HandleChannelChatMessage                  EventHandler[bindings.EventChannelChatMessage]                  `eventsub-type:"channel.chat.message" eventsub-version:"1"`
	HandleChannelChatClear                    EventHandler[bindings.EventChannelChatClear]                    `eventsub-type:"channel.chat.clear" eventsub-version:"1"`
	HandleChannelChatClearUserMessages        EventHandler[bindings.EventChannelChatClearUserMessages]        `eventsub-type:"channel.chat.clear_user_messages" eventsub-version:"1"`
	HandleChannelChatMessageDelete            EventHandler[bindings.EventChannelChatMessageDelete]            `eventsub-type:"channel.chat.message_delete" eventsub-version:"1"`
	HandleChannelChatNotification             EventHandler[bindings.EventChannelChatNotification]             `eventsub-type:"channel.chat.notification" eventsub-version:"1"`

	// ======================================================
	// NOTE: Beta handlers, may break backwards-compatibility
	// ======================================================
	HandleChannelUnbanRequestCreate  EventHandler[beta.EventChannelUnbanRequestCreate]  `eventsub-type:"channel.unban_request.create" eventsub-version:"beta"`
	HandleChannelUnbanRequestResolve EventHandler[beta.EventChannelUnbanRequestResolve] `eventsub-type:"channel.unban_request.resolve" eventsub-version:"beta"`
}

// NewHandler creates a new, empty EventSub handler. It implements http.Handler.
//
// When passed a non-nil secret, validation of notification events will be enabled
// and computed using the secret.
func NewHandler(secret []byte) *Handler {
	return &Handler{
		signatureSecret: secret,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.handlePost(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	// Read body into buffer
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = r.Body.Close()

	// Verify request signature
	if valid, err := h.verifySignature(r, bodyBytes); err != nil || !valid {
		http.Error(w, "Invalid request signature", http.StatusForbidden)
		return
	}

	// Decode request headers to verify and dispatch payload
	var headers bindings.NotificationHeaders
	if err := httpheader.Decode(r.Header, &headers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isDuplicate, err := h.checkIfDuplicate(r, &headers)
	if err != nil {
		// Error occurred while checking IDTracker
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if isDuplicate {
		// Call OnDuplicateNotification handler if set
		if h.OnDuplicateNotification != nil {
			h.OnDuplicateNotification(r.Context(), &headers)
		}
		writeEmptyOK(w) // ignore and return 2XX code
		return
	}

	switch headers.MessageType {
	case webhookCallbackVerification:
		h.handleVerification(w, r, bodyBytes, &headers)
		return
	case notificationMessageType:
		h.handleNotification(w, bodyBytes, &headers)
		return
	default:
		http.Error(w, "Unknown message type", http.StatusBadRequest)
		return
	}
}

// verifySignature verifies the Twitch-Eventsub-Message-Signature of the request.
// Returns whether the verification succeeded.
func (h *Handler) verifySignature(r *http.Request, body []byte) (bool, error) {
	if h.signatureSecret == nil {
		return true, nil
	}
	return VerifyRequestSignature(r, body, h.signatureSecret)
}

// checkIfDuplicate returns whether the IDTracker reports this notification is
// a duplicate. If it is a duplicate, it writes a 2xx response and returns true.
// Otherwise, it returns false.
func (h *Handler) checkIfDuplicate(r *http.Request, headers *bindings.NotificationHeaders) (bool, error) {
	if h.IDTracker == nil {
		return false, nil
	}

	duplicate, err := h.IDTracker.AddAndCheckIfDuplicate(r.Context(), headers.MessageID)
	if err != nil {
		return false, err
	}
	return duplicate, nil
}

func (h *Handler) handleVerification(
	w http.ResponseWriter,
	r *http.Request,
	bodyBytes []byte,
	headers *bindings.NotificationHeaders,
) {
	var data bindings.SubscriptionChallenge
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if h.VerifyChallenge == nil || h.VerifyChallenge(r.Context(), headers, &data) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(data.Challenge))
	} else {
		http.Error(w, "Invalid subscription", http.StatusBadRequest)
	}
}

// Writes a 200 OK response
func writeEmptyOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

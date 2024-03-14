package eventsub

//go:generate go run ./scripts/handler_generator --input=$GOFILE --output=notification_handler.go

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mozillazg/go-httpheader"

	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
)

const (
	webhookCallbackVerification = "webhook_callback_verification"
	notificationMessageType     = "notification"
)

type EventHandler[EventMessage any] func(h *bindings.NotificationHeaders, event *EventMessage)

// SubHandler implements http.Handler to receive Twitch webhook notifications.
//
// SubHandler handles both verification of new subscriptions and dispatching of
// event notifications. To handle a specific event, set the corresponding
// HandleXXX struct field. When a notification is received and validated, the
// handler function will be invoked in a new goroutine.
type SubHandler struct {
	doSignatureVerification bool
	signatureSecret         []byte

	// Challenge handler function.
	// Returns whether the subscription should be accepted.
	VerifyChallenge func(h *bindings.NotificationHeaders, chal *bindings.SubscriptionChallenge) bool

	// IDTracker used to deduplicate notifications
	IDTracker               IDTracker
	OnDuplicateNotification func(h *bindings.NotificationHeaders)

	HandleChannelUpdate                       EventHandler[bindings.EventChannelUpdate]                       `eventsub-type:"channel.update"`
	HandleChannelFollow                       EventHandler[bindings.EventChannelFollow]                       `eventsub-type:"channel.follow"`
	HandleChannelSubscribe                    EventHandler[bindings.EventChannelSubscribe]                    `eventsub-type:"channel.subscribe"`
	HandleChannelSubscriptionEnd              EventHandler[bindings.EventChannelSubscriptionEnd]              `eventsub-type:"channel.subscription.end"`
	HandleChannelSubscriptionGift             EventHandler[bindings.EventChannelSubscriptionGift]             `eventsub-type:"channel.subscription.gift"`
	HandleChannelSubscriptionMessage          EventHandler[bindings.EventChannelSubscriptionMessage]          `eventsub-type:"channel.subscription.message"`
	HandleChannelCheer                        EventHandler[bindings.EventChannelCheer]                        `eventsub-type:"channel.cheer"`
	HandleChannelRaid                         EventHandler[bindings.EventChannelRaid]                         `eventsub-type:"channel.raid"`
	HandleChannelBan                          EventHandler[bindings.EventChannelBan]                          `eventsub-type:"channel.ban"`
	HandleChannelUnban                        EventHandler[bindings.EventChannelUnban]                        `eventsub-type:"channel.unban"`
	HandleChannelUnbanRequestCreate           EventHandler[bindings.EventChannelUnbanRequestCreate]           `eventsub-type:"channel.unban_request.create"`
	HandleChannelUnbanRequestResolve          EventHandler[bindings.EventChannelUnbanRequestResolve]          `eventsub-type:"channel.unban_request.resolve"`
	HandleChannelModeratorAdd                 EventHandler[bindings.EventChannelModeratorAdd]                 `eventsub-type:"channel.moderator.add"`
	HandleChannelModeratorRemove              EventHandler[bindings.EventChannelModeratorRemove]              `eventsub-type:"channel.moderator.remove"`
	HandleChannelPointsRewardAdd              EventHandler[bindings.EventChannelPointsRewardAdd]              `eventsub-type:"channel.channel_points_custom_reward.add"`
	HandleChannelPointsRewardUpdate           EventHandler[bindings.EventChannelPointsRewardUpdate]           `eventsub-type:"channel.channel_points_custom_reward.update"`
	HandleChannelPointsRewardRemove           EventHandler[bindings.EventChannelPointsRewardRemove]           `eventsub-type:"channel.channel_points_custom_reward.remove"`
	HandleChannelPointsRewardRedemptionAdd    EventHandler[bindings.EventChannelPointsRewardRedemptionAdd]    `eventsub-type:"channel.channel_points_custom_reward_redemption.add"`
	HandleChannelPointsRewardRedemptionUpdate EventHandler[bindings.EventChannelPointsRewardRedemptionUpdate] `eventsub-type:"channel.channel_points_custom_reward_redemption.update"`
	HandleChannelPollBegin                    EventHandler[bindings.EventChannelPollBegin]                    `eventsub-type:"channel.poll.begin"`
	HandleChannelPollProgress                 EventHandler[bindings.EventChannelPollProgress]                 `eventsub-type:"channel.poll.progress"`
	HandleChannelPollEnd                      EventHandler[bindings.EventChannelPollEnd]                      `eventsub-type:"channel.poll.end"`
	HandleChannelPredictionBegin              EventHandler[bindings.EventChannelPredictionBegin]              `eventsub-type:"channel.prediction.begin"`
	HandleChannelPredictionProgress           EventHandler[bindings.EventChannelPredictionProgress]           `eventsub-type:"channel.prediction.progress"`
	HandleChannelPredictionLock               EventHandler[bindings.EventChannelPredictionLock]               `eventsub-type:"channel.prediction.lock"`
	HandleChannelPredictionEnd                EventHandler[bindings.EventChannelPredictionEnd]                `eventsub-type:"channel.prediction.end"`
	HandleDropEntitlementGrant                EventHandler[bindings.EventDropEntitlementGrant]                `eventsub-type:"drop.entitlement.grant"`
	HandleExtensionBitsTransactionCreate      EventHandler[bindings.EventBitsTransactionCreate]               `eventsub-type:"extension.bits_transaction.create"`
	HandleGoalBegin                           EventHandler[bindings.EventGoals]                               `eventsub-type:"channel.goal.begin"`
	HandleGoalProgress                        EventHandler[bindings.EventGoals]                               `eventsub-type:"channel.goal.progress"`
	HandleGoalEnd                             EventHandler[bindings.EventGoals]                               `eventsub-type:"channel.goal.end"`
	HandleHypeTrainBegin                      EventHandler[bindings.EventHypeTrainBegin]                      `eventsub-type:"channel.hype_train.begin"`
	HandleHypeTrainProgress                   EventHandler[bindings.EventHypeTrainProgress]                   `eventsub-type:"channel.hype_train.progress"`
	HandleHypeTrainEnd                        EventHandler[bindings.EventHypeTrainEnd]                        `eventsub-type:"channel.hype_train.end"`
	HandleStreamOnline                        EventHandler[bindings.EventStreamOnline]                        `eventsub-type:"stream.online"`
	HandleStreamOffline                       EventHandler[bindings.EventStreamOffline]                       `eventsub-type:"stream.offline"`
	HandleUserUpdate                          EventHandler[bindings.EventUserUpdate]                          `eventsub-type:"user.update"`
	HandleUserAuthorizationGrant              EventHandler[bindings.EventUserAuthorizationGrant]              `eventsub-type:"user.authorization.grant"`
	HandleUserAuthorizationRevoke             EventHandler[bindings.EventUserAuthorizationRevoke]             `eventsub-type:"user.authorization.revoke"`
	HandleChannelChatMessage                  EventHandler[bindings.EventChannelChatMessage]                  `eventsub-type:"channel.chat.message"`
	HandleChannelChatClear                    EventHandler[bindings.EventChannelChatClear]                    `eventsub-type:"channel.chat.clear"`
	HandleChannelChatClearUserMessages        EventHandler[bindings.EventChannelChatClearUserMessages]        `eventsub-type:"channel.chat.clear_user_messages"`
	HandleChannelChatMessageDelete            EventHandler[bindings.EventChannelChatMessageDelete]            `eventsub-type:"channel.chat.message_delete"`
	HandleChannelChatNotification             EventHandler[bindings.EventChannelChatNotification]             `eventsub-type:"channel.chat.notification"`
}

func NewSubHandler(doSignatureVerification bool, secret []byte) *SubHandler {
	if doSignatureVerification && secret == nil {
		panic("secret must be set if signature verification is enabled")
	}

	return &SubHandler{
		doSignatureVerification: doSignatureVerification,
		signatureSecret:         secret,
	}
}

func (s *SubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		s.handlePost(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *SubHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	// Read body into buffer
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if s.doSignatureVerification {
		if valid, err := VerifyRequestSignature(r, bodyBytes, s.signatureSecret); err != nil || !valid {
			http.Error(w, "Invalid request signature", http.StatusForbidden)
			return
		}
	}

	// Decode request headers to verify and dispatch payload
	var h bindings.NotificationHeaders
	if err := httpheader.Decode(r.Header, &h); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isDuplicate, err := s.checkIfDuplicate(w, r, &h)
	if err != nil {
		// Error occurred while checking IDTracker
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if isDuplicate {
		return // already handled response
	}

	switch h.MessageType {
	case webhookCallbackVerification:
		s.handleVerification(w, bodyBytes, &h)
		return
	case notificationMessageType:
		s.handleNotification(w, bodyBytes, &h)
		return
	default:
		http.Error(w, "Unknown message type", http.StatusBadRequest)
		return
	}
}

// checkIfDuplicate returns whether the IDTracker reports this notification is
// a duplicate. If it is a duplicate, it writes a 2xx response and returns true.
// Otherwise, it returns false.
func (s *SubHandler) checkIfDuplicate(
	w http.ResponseWriter,
	r *http.Request,
	h *bindings.NotificationHeaders,
) (bool, error) {
	if s.IDTracker != nil {
		duplicate, err := s.IDTracker.AddAndCheckIfDuplicate(r.Context(), h.MessageID)
		if err != nil {
			return false, err
		}

		if duplicate {
			if s.OnDuplicateNotification != nil {
				go s.OnDuplicateNotification(h)
			}
			writeEmptyOK(w) // ignore and return 2XX code
			return true, nil
		}
	}

	return false, nil
}

func (s *SubHandler) handleVerification(
	w http.ResponseWriter,
	bodyBytes []byte,
	headers *bindings.NotificationHeaders,
) {
	var data bindings.SubscriptionChallenge
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if s.VerifyChallenge == nil || s.VerifyChallenge(headers, &data) {
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

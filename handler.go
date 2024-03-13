package eventsub

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
	VerifyChallenge func(h *bindings.ResponseHeaders, chal *bindings.SubscriptionChallenge) bool

	// IDTracker used to deduplicate notifications
	IDTracker               IDTracker
	OnDuplicateNotification func(h *bindings.ResponseHeaders)

	HandleChannelUpdate func(h *bindings.ResponseHeaders, event *bindings.EventChannelUpdate)
	HandleChannelFollow func(h *bindings.ResponseHeaders, event *bindings.EventChannelFollow)
	HandleUserUpdate    func(h *bindings.ResponseHeaders, event *bindings.EventUserUpdate)

	HandleChannelSubscribe       func(h *bindings.ResponseHeaders, event *bindings.EventChannelSubscribe)
	HandleChannelSubscriptionEnd func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelSubscriptionEnd,
	)
	HandleChannelSubscriptionGift func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelSubscriptionGift,
	)
	HandleChannelSubscriptionMessage func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelSubscriptionMessage,
	)
	HandleChannelCheer func(h *bindings.ResponseHeaders, event *bindings.EventChannelCheer)
	HandleChannelRaid  func(h *bindings.ResponseHeaders, event *bindings.EventChannelRaid)

	HandleChannelBan             func(h *bindings.ResponseHeaders, event *bindings.EventChannelBan)
	HandleChannelUnban           func(h *bindings.ResponseHeaders, event *bindings.EventChannelUnban)
	HandleChannelModeratorAdd    func(h *bindings.ResponseHeaders, event *bindings.EventChannelModeratorAdd)
	HandleChannelModeratorRemove func(h *bindings.ResponseHeaders, event *bindings.EventChannelModeratorRemove)

	HandleChannelPointsRewardAdd func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPointsRewardAdd,
	)
	HandleChannelPointsRewardUpdate func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPointsRewardUpdate,
	)
	HandleChannelPointsRewardRemove func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPointsRewardRemove,
	)
	HandleChannelPointsRewardRedemptionAdd func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPointsRewardRedemptionAdd,
	)
	HandleChannelPointsRewardRedemptionUpdate func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPointsRewardRedemptionUpdate,
	)

	HandleChannelPollBegin    func(h *bindings.ResponseHeaders, event *bindings.EventChannelPollBegin)
	HandleChannelPollProgress func(h *bindings.ResponseHeaders, event *bindings.EventChannelPollProgress)
	HandleChannelPollEnd      func(h *bindings.ResponseHeaders, event *bindings.EventChannelPollEnd)

	HandleChannelPredictionBegin func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPredictionBegin,
	)
	HandleChannelPredictionProgress func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPredictionProgress,
	)
	HandleChannelPredictionLock func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelPredictionLock,
	)
	HandleChannelPredictionEnd func(h *bindings.ResponseHeaders, event *bindings.EventChannelPredictionEnd)

	HandleDropEntitlementGrant func(
		h *bindings.ResponseHeaders,
		event *bindings.EventDropEntitlementGrant,
	)
	HandleExtensionBitsTransactionCreate func(
		h *bindings.ResponseHeaders,
		event *bindings.EventBitsTransactionCreate,
	)

	HandleGoalBegin    func(h *bindings.ResponseHeaders, event *bindings.EventGoals)
	HandleGoalProgress func(h *bindings.ResponseHeaders, event *bindings.EventGoals)
	HandleGoalEnd      func(h *bindings.ResponseHeaders, event *bindings.EventGoals)

	HandleHypeTrainBegin    func(h *bindings.ResponseHeaders, event *bindings.EventHypeTrainBegin)
	HandleHypeTrainProgress func(h *bindings.ResponseHeaders, event *bindings.EventHypeTrainProgress)
	HandleHypeTrainEnd      func(h *bindings.ResponseHeaders, event *bindings.EventHypeTrainEnd)

	HandleStreamOnline  func(h *bindings.ResponseHeaders, event *bindings.EventStreamOnline)
	HandleStreamOffline func(h *bindings.ResponseHeaders, event *bindings.EventStreamOffline)

	HandleUserAuthorizationGrant  func(h *bindings.ResponseHeaders, event *bindings.EventUserAuthorizationGrant)
	HandleUserAuthorizationRevoke func(
		h *bindings.ResponseHeaders,
		event *bindings.EventUserAuthorizationRevoke,
	)

	HandleChannelChatMessage           func(h *bindings.ResponseHeaders, event *bindings.EventChannelChatMessage)
	HandleChannelChatClear             func(h *bindings.ResponseHeaders, event *bindings.EventChannelChatClear)
	HandleChannelChatClearUserMessages func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelChatClearUserMessages,
	)
	HandleChannelChatMessageDelete func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelChatMessageDelete,
	)

	HandleChannelChatNotification func(
		h *bindings.ResponseHeaders,
		event *bindings.EventChannelChatNotification,
	)
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
	var h bindings.ResponseHeaders
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
	h *bindings.ResponseHeaders,
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
	headers *bindings.ResponseHeaders,
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

func (s *SubHandler) handleNotification(
	w http.ResponseWriter,
	bodyBytes []byte,
	h *bindings.ResponseHeaders,
) {
	var notification bindings.EventNotification
	if err := json.Unmarshal(bodyBytes, &notification); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	event := notification.Event

	switch h.SubscriptionType {
	case "channel.update":
		var data bindings.EventChannelUpdate
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelUpdate != nil {
			go s.HandleChannelUpdate(h, &data)
		}
	case "channel.follow":
		var data bindings.EventChannelFollow
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelFollow != nil {
			go s.HandleChannelFollow(h, &data)
		}
	case "channel.subscribe":
		var data bindings.EventChannelSubscribe
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelSubscribe != nil {
			go s.HandleChannelSubscribe(h, &data)
		}
	case "channel.subscription.end":
		var data bindings.EventChannelSubscriptionEnd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelSubscriptionEnd != nil {
			go s.HandleChannelSubscriptionEnd(h, &data)
		}
	case "channel.subscription.gift":
		var data bindings.EventChannelSubscriptionGift
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelSubscriptionGift != nil {
			go s.HandleChannelSubscriptionGift(h, &data)
		}
	case "channel.subscription.message":
		var data bindings.EventChannelSubscriptionMessage
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelSubscriptionMessage != nil {
			go s.HandleChannelSubscriptionMessage(h, &data)
		}
	case "channel.cheer":
		var data bindings.EventChannelCheer
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelCheer != nil {
			go s.HandleChannelCheer(h, &data)
		}
	case "channel.raid":
		var data bindings.EventChannelRaid
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelRaid != nil {
			go s.HandleChannelRaid(h, &data)
		}
	case "channel.ban":
		var data bindings.EventChannelBan
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelBan != nil {
			go s.HandleChannelBan(h, &data)
		}
	case "channel.unban":
		var data bindings.EventChannelUnban
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelUnban != nil {
			go s.HandleChannelUnban(h, &data)
		}
	case "channel.moderator.add":
		var data bindings.EventChannelModeratorAdd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelModeratorAdd != nil {
			go s.HandleChannelModeratorAdd(h, &data)
		}
	case "channel.moderator.remove":
		var data bindings.EventChannelModeratorRemove
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelModeratorRemove != nil {
			go s.HandleChannelModeratorRemove(h, &data)
		}
	case "channel.channel_points_custom_reward.add":
		var data bindings.EventChannelPointsRewardAdd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPointsRewardAdd != nil {
			go s.HandleChannelPointsRewardAdd(h, &data)
		}
	case "channel.channel_points_custom_reward.update":
		var data bindings.EventChannelPointsRewardUpdate
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPointsRewardUpdate != nil {
			go s.HandleChannelPointsRewardUpdate(h, &data)
		}
	case "channel.channel_points_custom_reward.remove":
		var data bindings.EventChannelPointsRewardRemove
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPointsRewardRemove != nil {
			go s.HandleChannelPointsRewardRemove(h, &data)
		}
	case "channel.channel_points_custom_reward_redemption.add":
		var data bindings.EventChannelPointsRewardRedemptionAdd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPointsRewardRedemptionAdd != nil {
			go s.HandleChannelPointsRewardRedemptionAdd(h, &data)
		}
	case "channel.channel_points_custom_reward_redemption.update":
		var data bindings.EventChannelPointsRewardRedemptionUpdate
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPointsRewardRedemptionUpdate != nil {
			go s.HandleChannelPointsRewardRedemptionUpdate(h, &data)
		}
	case "channel.poll.begin":
		var data bindings.EventChannelPollBegin
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPollBegin != nil {
			go s.HandleChannelPollBegin(h, &data)
		}
	case "channel.poll.progress":
		var data bindings.EventChannelPollProgress
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPollProgress != nil {
			go s.HandleChannelPollProgress(h, &data)
		}
	case "channel.poll.end":
		var data bindings.EventChannelPollEnd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPollEnd != nil {
			go s.HandleChannelPollEnd(h, &data)
		}
	case "channel.prediction.begin":
		var data bindings.EventChannelPredictionBegin
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPredictionBegin != nil {
			go s.HandleChannelPredictionBegin(h, &data)
		}
	case "channel.prediction.progress":
		var data bindings.EventChannelPredictionProgress
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPredictionProgress != nil {
			go s.HandleChannelPredictionProgress(h, &data)
		}
	case "channel.prediction.lock":
		var data bindings.EventChannelPredictionLock
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPredictionLock != nil {
			go s.HandleChannelPredictionLock(h, &data)
		}
	case "channel.prediction.end":
		var data bindings.EventChannelPredictionEnd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelPredictionEnd != nil {
			go s.HandleChannelPredictionEnd(h, &data)
		}
	case "drop.entitlement.grant":
		var data bindings.EventDropEntitlementGrant
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleDropEntitlementGrant != nil {
			go s.HandleDropEntitlementGrant(h, &data)
		}
	case "extension.bits_transaction.create":
		var data bindings.EventBitsTransactionCreate
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleExtensionBitsTransactionCreate != nil {
			go s.HandleExtensionBitsTransactionCreate(h, &data)
		}
	case "channel.goal.begin":
		var data bindings.EventGoals
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleGoalBegin != nil {
			go s.HandleGoalBegin(h, &data)
		}
	case "channel.goal.progress":
		var data bindings.EventGoals
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleGoalProgress != nil {
			go s.HandleGoalProgress(h, &data)
		}
	case "channel.goal.end":
		var data bindings.EventGoals
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleGoalEnd != nil {
			go s.HandleGoalEnd(h, &data)
		}
	case "channel.hype_train.begin":
		var data bindings.EventHypeTrainBegin
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleHypeTrainBegin != nil {
			go s.HandleHypeTrainBegin(h, &data)
		}
	case "channel.hype_train.progress":
		var data bindings.EventHypeTrainProgress
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleHypeTrainProgress != nil {
			go s.HandleHypeTrainProgress(h, &data)
		}
	case "channel.hype_train.end":
		var data bindings.EventHypeTrainEnd
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleHypeTrainEnd != nil {
			go s.HandleHypeTrainEnd(h, &data)
		}
	case "stream.online":
		var data bindings.EventStreamOnline
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleStreamOnline != nil {
			go s.HandleStreamOnline(h, &data)
		}
	case "stream.offline":
		var data bindings.EventStreamOffline
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleStreamOffline != nil {
			go s.HandleStreamOffline(h, &data)
		}
	case "user.authorization.grant":
		var data bindings.EventUserAuthorizationGrant
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleUserAuthorizationGrant != nil {
			go s.HandleUserAuthorizationGrant(h, &data)
		}
	case "user.authorization.revoke":
		var data bindings.EventUserAuthorizationRevoke
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleUserAuthorizationRevoke != nil {
			go s.HandleUserAuthorizationRevoke(h, &data)
		}
	case "user.update":
		var data bindings.EventUserUpdate
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleUserUpdate != nil {
			go s.HandleUserUpdate(h, &data)
		}
	case "channel.chat.message":
		var data bindings.EventChannelChatMessage
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelChatMessage != nil {
			go s.HandleChannelChatMessage(h, &data)
		}
	case "channel.chat.clear":
		var data bindings.EventChannelChatClear
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelChatClear != nil {
			go s.HandleChannelChatClear(h, &data)
		}
	case "channel.chat.clear_user_messages":
		var data bindings.EventChannelChatClearUserMessages
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelChatClearUserMessages != nil {
			go s.HandleChannelChatClearUserMessages(h, &data)
		}
	case "channel.chat.message_delete":
		var data bindings.EventChannelChatMessageDelete
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelChatMessageDelete != nil {
			go s.HandleChannelChatMessageDelete(h, &data)
		}
	case "channel.chat.notification":
		var data bindings.EventChannelChatNotification
		if err := json.Unmarshal(event, &data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		if s.HandleChannelChatNotification != nil {
			go s.HandleChannelChatNotification(h, &data)
		}
	default:
		http.Error(w, "Unknown notification type", http.StatusBadRequest)
		return
	}

	writeEmptyOK(w)
}

// Writes a 200 OK response
func writeEmptyOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

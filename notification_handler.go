// Code generated by handler_generator. DO NOT EDIT.

package eventsub

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
)

func deserializeAndCallHandler[EventType any](
	h *bindings.NotificationHeaders,
	notification bindings.EventNotification,
	handler EventHandler[EventType],
) error {
	if handler == nil {
		return nil
	}

	var data EventType
	if err := json.Unmarshal(notification.Event, &data); err != nil {
		return err
	}

	go handler(*h, notification.Subscription, data)
	return nil
}

func (h *Handler) handleNotification(
	ctx context.Context,
	w http.ResponseWriter,
	bodyBytes []byte,
	headers *bindings.NotificationHeaders,
) {
	var notification bindings.EventNotification
	if err := json.Unmarshal(bodyBytes, &notification); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if h.BeforeHandleEvent != nil {
		h.BeforeHandleEvent(ctx, headers, &notification)
	}

	var err error
	selector := headers.SubscriptionType + "_" + headers.SubscriptionVersion
	switch selector {
	case "channel.ban_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelBan)

	case "channel.channel_points_custom_reward.add_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPointsRewardAdd)

	case "channel.channel_points_custom_reward.remove_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPointsRewardRemove)

	case "channel.channel_points_custom_reward.update_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPointsRewardUpdate)

	case "channel.channel_points_custom_reward_redemption.add_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPointsRewardRedemptionAdd)

	case "channel.channel_points_custom_reward_redemption.update_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPointsRewardRedemptionUpdate)

	case "channel.chat.clear_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelChatClear)

	case "channel.chat.clear_user_messages_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelChatClearUserMessages)

	case "channel.chat.message_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelChatMessage)

	case "channel.chat.message_delete_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelChatMessageDelete)

	case "channel.chat.notification_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelChatNotification)

	case "channel.cheer_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelCheer)

	case "channel.follow_2":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelFollow)

	case "channel.goal.begin_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleGoalBegin)

	case "channel.goal.end_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleGoalEnd)

	case "channel.goal.progress_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleGoalProgress)

	case "channel.hype_train.begin_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleHypeTrainBegin)

	case "channel.hype_train.end_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleHypeTrainEnd)

	case "channel.hype_train.progress_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleHypeTrainProgress)

	case "channel.moderator.add_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelModeratorAdd)

	case "channel.moderator.remove_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelModeratorRemove)

	case "channel.poll.begin_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPollBegin)

	case "channel.poll.end_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPollEnd)

	case "channel.poll.progress_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPollProgress)

	case "channel.prediction.begin_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPredictionBegin)

	case "channel.prediction.end_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPredictionEnd)

	case "channel.prediction.lock_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPredictionLock)

	case "channel.prediction.progress_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelPredictionProgress)

	case "channel.raid_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelRaid)

	case "channel.subscribe_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelSubscribe)

	case "channel.subscription.end_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelSubscriptionEnd)

	case "channel.subscription.gift_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelSubscriptionGift)

	case "channel.subscription.message_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelSubscriptionMessage)

	case "channel.unban_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelUnban)

	case "channel.unban_request.create_beta":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelUnbanRequestCreate)

	case "channel.unban_request.resolve_beta":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelUnbanRequestResolve)

	case "channel.update_2":
		err = deserializeAndCallHandler(headers, notification, h.HandleChannelUpdate)

	case "drop.entitlement.grant_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleDropEntitlementGrant)

	case "extension.bits_transaction.create_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleExtensionBitsTransactionCreate)

	case "stream.offline_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleStreamOffline)

	case "stream.online_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleStreamOnline)

	case "user.authorization.grant_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleUserAuthorizationGrant)

	case "user.authorization.revoke_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleUserAuthorizationRevoke)

	case "user.update_1":
		err = deserializeAndCallHandler(headers, notification, h.HandleUserUpdate)

	default:
		http.Error(w, "Unsupported notification type and version", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Invalid notification", http.StatusBadRequest)
		return
	}

	writeEmptyOK(w)
}

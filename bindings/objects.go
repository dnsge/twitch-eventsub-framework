package bindings

type PollChoice struct {
	// ID for the choice.
	ID string `json:"id"`
	// Text displayed for the choice.
	Title string `json:"title"`
	// Number of votes received via Bits.
	BitsVotes int `json:"bits_votes"`
	// Number of votes received via Channel Points.
	ChannelPointsVotes int `json:"channel_points_votes"`
	// Total number of votes received for the choice across all methods of voting.
	Votes int `json:"votes"`
}

type BitsVoting struct {
	// Indicates if Bits can be used for voting.
	IsEnabled bool `json:"is_enabled"`
	// Number of Bits required to vote once with Bits.
	AmountPerVote int `json:"amount_per_vote"`
}

type ChannelPointsVoting struct {
	// Indicates if Channel Points can be used for voting.
	IsEnabled bool `json:"is_enabled"`
	// Number of Channel Points required to vote once with Channel Points.
	AmountPerVote int `json:"amount_per_vote"`
}

type MaxPerStream struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The max per stream limit.
	Value int `json:"value"`
}

type MaxPerUserPerStream struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The max per user per stream limit.
	Value int `json:"value"`
}

type Image struct {
	// URL for the image at 1x size.
	Url1x string `json:"url_1x"`
	// URL for the image at 2x size.
	Url2x string `json:"url_2x"`
	// URL for the image at 4x size.
	Url4x string `json:"url_4x"`
}

type GlobalCooldown struct {
	// Is the setting enabled.
	IsEnabled bool `json:"is_enabled"`
	// The cooldown in seconds.
	Seconds int `json:"seconds"`
}

type Reward struct {
	// The reward identifier.
	ID string `json:"id"`
	// The reward name.
	Title string `json:"title"`
	// The reward cost.
	Cost int `json:"cost"`
	// The reward description.
	Prompt string `json:"prompt"`
}

type PredictionOutcome struct {
	// The outcome ID.
	ID string `json:"id"`
	// The outcome title.
	Title string `json:"title"`
	// The color for the outcome. Valid values are pink and blue.
	Color string `json:"color"`
	// The number of users who used Channel Points on this outcome.
	Users int `json:"users"`
	// The total number of Channel Points used on this outcome.
	ChannelPoints int `json:"channel_points"`
	// An array of users who used the most Channel Points on this outcome.
	TopPredictors []TopPredictor `json:"top_predictors"`
}

type TopPredictor struct {
	// The ID of the user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// The number of Channel Points won.
	// This value is always null in the event payload for Prediction progress and Prediction lock.
	// This value is 0 if the outcome did not win or if the Prediction was canceled and Channel Points were refunded.
	ChannelPointsWon int `json:"channel_points_won"`
	// The number of Channel Points used to participate in the Prediction.
	ChannelPointsUsed int `json:"channel_points_used"`
}

type Message struct {
	// The text of the resubscription chat message.
	Text string `json:"text"`
	// An array that includes the emote ID and start and end positions for where the emote appears in the text.
	Emotes []Emote `json:"emotes"`
}

type Emote struct {
	// The index of where the ChatNotificationMessageFragmentEmote starts in the text.
	Begin int `json:"begin"`
	// The index of where the ChatNotificationMessageFragmentEmote ends in the text.
	End int `json:"end"`
	// The emote ID.
	ID string `json:"id"`
}

type Product struct {
	// Product name.
	Name string `json:"name"`
	// Bits involved in the transaction.
	Bits int `json:"bits"`
	// Unique identifier for the product acquired.
	Sku string `json:"sku"`
	// Flag indicating if the product is in development. If InDevelopment is true, bits will be 0.
	InDevelopment bool `json:"in_development"`
}

type TopContributor struct {
	// The ID of the user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// Type of contribution. Valid values include bits, subscription.
	Type string `json:"type"`
	// The total contributed.
	Total int `json:"total"`
}

type LastContribution struct {
	// The ID of the user.
	UserID string `json:"user_id"`
	// The login of the user.
	UserLogin string `json:"user_login"`
	// The display name of the user.
	UserName string `json:"user_name"`
	// Type of contribution. Valid values include bits, subscription.
	Type string `json:"type"`
	// The total contributed.
	Total int `json:"total"`
}

// ChatNotificationBadge ChatBadge represents a chat badge.
type ChatNotificationBadge struct {
	// An ID that identifies this set of chat badges. For example, Bits or Subscriber.
	SetID string `json:"set_id"`
	// An ID that identifies this version of the badge. The ID can be any value. For example, for Bits, the ID is the Bits tier level, but for World of Warcraft, it could be Alliance or Horde.
	ID string `json:"id"`
	// Contains metadata related to the chat badges in the badges tag. Currently, this tag contains metadata only for subscriber badges, to indicate the number of months the user has been a subscriber.
	Info string `json:"info"`
}

// ChatNotificationMessage ChatMessage represents a structured chat message.
type ChatNotificationMessage struct {
	// The chat message in plain text.
	Text string `json:"text"`
	// Ordered list of chat message fragments.
	Fragments []ChatNotificationMessageFragment `json:"fragments"`
}

// ChatNotificationMessageFragment represents a fragment of a chat message.
type ChatNotificationMessageFragment struct {
	// The type of message fragment. Possible values:
	//  - text
	//  - cheermote
	//  - emote
	//  - mention
	Type string `json:"type"`
	// Message text in fragment
	Text string `json:"text"`
	// Optional. Metadata pertaining to the cheermote.
	Cheermote *ChatNotificationMessageFragmentCheermote `json:"cheermote,omitempty"`
	// Optional. Metadata pertaining to the emote.
	Emote *ChatNotificationMessageFragmentEmote `json:"emote,omitempty"`
	// Optional. Metadata pertaining to the mention.
	Mention *ChatNotificationMessageFragmentMention `json:"mention,omitempty"`
}

// ChatNotificationMessageFragmentCheermote represents metadata pertaining to a cheermote.
type ChatNotificationMessageFragmentCheermote struct {
	// The name portion of the ChatNotificationMessageFragmentCheermote string that you use in chat to cheer Bits. The full ChatNotificationMessageFragmentCheermote string is the concatenation of {prefix} + {number of Bits}. For example, if the prefix is “Cheer” and you want to cheer 100 Bits, the full ChatNotificationMessageFragmentCheermote string is Cheer100. When the ChatNotificationMessageFragmentCheermote string is entered in chat, Twitch converts it to the image associated with the Bits tier that was cheered.
	Prefix string `json:"prefix"`
	// The amount of bits cheered.
	Bits int `json:"bits"`
	// The tier level of the cheermote.
	Tier int `json:"tier"`
}

// ChatNotificationMessageFragmentEmote represents metadata pertaining to an emote.
type ChatNotificationMessageFragmentEmote struct {
	// An ID that uniquely identifies this emote.
	ID string `json:"id"`
	// An ID that identifies the emote set that the emote belongs to.
	EmoteSetID string `json:"emote_set_id"`
	// The ID of the broadcaster who owns the emote.
	OwnerID string `json:"owner_id"`
	// The formats that the emote is available in. For example, if the emote is available only as a static PNG, the array contains only static. But if the emote is available as a static PNG and an animated GIF, the array contains static and animated. The possible formats are:
	//  - animated — An animated GIF is available for this emote.
	//  - static — A static PNG file is available for this emote.
	Format []string `json:"format"`
}

// ChatNotificationMessageFragmentMention represents metadata pertaining to a mention.
type ChatNotificationMessageFragmentMention struct {
	// The user ID of the mentioned user.
	UserID string `json:"user_id"`
	// The user name of the mentioned user.
	UserName string `json:"user_name"`
	// The user login of the mentioned user.
	UserLogin string `json:"user_login"`
}

// ChatNotificationSubEvent SubEvent represents information about the sub event.
type ChatNotificationSubEvent struct {
	// The type of subscription plan being used. Possible values are:
	//  - 1000 — First level of paid or Prime subscription
	//  - 2000 — Second level of paid subscription
	//  - 3000 — Third level of paid subscription
	SubTier string `json:"sub_tier"`

	// Indicates if the subscription was obtained through Amazon Prime.
	IsPrime bool `json:"is_prime"`

	// The number of months the subscription is for.
	DurationMonths int `json:"duration_months"`
}

// ChatNotificationResubEvent represents information about the resub event.
type ChatNotificationResubEvent struct {
	// The total number of months the user has subscribed.
	CumulativeMonths int `json:"cumulative_months"`

	// The number of months the subscription is for.
	DurationMonths int `json:"duration_months"`

	// Optional. The number of consecutive months the user has subscribed.
	StreakMonths int `json:"streak_months,omitempty"`

	// The type of subscription plan being used. Possible values are:
	//  - 1000 — First level of paid or Prime subscription
	//  - 2000 — Second level of paid subscription
	//  - 3000 — Third level of paid subscription
	SubTier string `json:"sub_tier"`

	// Indicates if the resub was obtained through Amazon Prime.
	IsPrime bool `json:"is_prime"`

	// Whether or not the resub was a result of a gift.
	IsGift bool `json:"is_gift"`

	// Optional. Whether or not the gift was anonymous.
	GifterIsAnonymous bool `json:"gifter_is_anonymous,omitempty"`

	// Optional. The user ID of the subscription gifter. Null if anonymous.
	GifterUserID string `json:"gifter_user_id,omitempty"`

	// Optional. The user name of the subscription gifter. Null if anonymous.
	GifterUserName string `json:"gifter_user_name,omitempty"`

	// Optional. The user login of the subscription gifter. Null if anonymous.
	GifterUserLogin string `json:"gifter_user_login,omitempty"`
}

// ChatNotificationSubGiftEvent represents information about the gift sub event.
type ChatNotificationSubGiftEvent struct {
	// The number of months the subscription is for.
	DurationMonths int `json:"duration_months"`

	// Optional. The amount of gifts the gifter has given in this channel. Null if anonymous.
	CumulativeTotal int `json:"cumulative_total,omitempty"`

	// The user ID of the subscription gift recipient.
	RecipientUserID string `json:"recipient_user_id"`

	// The user name of the subscription gift recipient.
	RecipientUserName string `json:"recipient_user_name"`

	// The user login of the subscription gift recipient.
	RecipientUserLogin string `json:"recipient_user_login"`

	// The type of subscription plan being used. Possible values are:
	//  - 1000 — First level of paid subscription
	//  - 2000 — Second level of paid subscription
	//  - 3000 — Third level of paid subscription
	SubTier string `json:"sub_tier"`

	// Optional. The ID of the associated community gift. Null if not associated with a community gift.
	CommunityGiftID string `json:"community_gift_id,omitempty"`
}

// ChatNotificationCommunitySubGiftEvent represents information about the community gift sub event.
type ChatNotificationCommunitySubGiftEvent struct {
	// The ID of the associated community gift.
	ID string `json:"id"`

	// Number of subscriptions being gifted.
	Total int `json:"total"`

	// The type of subscription plan being used. Possible values are:
	//  - 1000 — First level of paid subscription
	//  - 2000 — Second level of paid subscription
	//  - 3000 — Third level of paid subscription
	SubTier string `json:"sub_tier"`

	// Optional. The amount of gifts the gifter has given in this channel. Null if anonymous.
	CumulativeTotal int `json:"cumulative_total,omitempty"`
}

// ChatNotificationGiftPaidUpgradeEvent represents information about the community gift paid upgrade event.
type ChatNotificationGiftPaidUpgradeEvent struct {
	// Whether the gift was given anonymously.
	GifterIsAnonymous bool `json:"gifter_is_anonymous"`

	// Optional. The user ID of the user who gifted the subscription. Null if anonymous.
	GifterUserID string `json:"gifter_user_id,omitempty"`

	// Optional. The user name of the user who gifted the subscription. Null if anonymous.
	GifterUserName string `json:"gifter_user_name,omitempty"`

	// Optional. The user login of the user who gifted the subscription. Null if anonymous.
	GifterUserLogin string `json:"gifter_user_login,omitempty"`
}

// ChatNotificationPrimePaidUpgradeEvent represents information about the Prime gift paid upgrade event.
type ChatNotificationPrimePaidUpgradeEvent struct {
	// The type of subscription plan being used. Possible values are:
	//  - 1000 — First level of paid subscription
	//  - 2000 — Second level of paid subscription
	//  - 3000 — Third level of paid subscription
	SubTier string `json:"sub_tier"`
}

// ChatNotificationRaidEvent represents information about the raid event.
type ChatNotificationRaidEvent struct {
	// The user ID of the broadcaster raiding this channel.
	UserID string `json:"user_id"`

	// The user name of the broadcaster raiding this channel.
	UserName string `json:"user_name"`

	// The login name of the broadcaster raiding this channel.
	UserLogin string `json:"user_login"`

	// The number of viewers raiding this channel from the broadcaster’s channel.
	ViewerCount int `json:"viewer_count"`

	// Profile image URL of the broadcaster raiding this channel.
	ProfileImageURL string `json:"profile_image_url"`
}

// ChatNotificationUnraidEvent represents an empty payload for the unraid event.
type ChatNotificationUnraidEvent struct {
}

// ChatNotificationPayItForwardEvent represents information about the pay it forward event.
type ChatNotificationPayItForwardEvent struct {
	// Whether the gift was given anonymously.
	GifterIsAnonymous bool `json:"gifter_is_anonymous"`

	// Optional. The user ID of the user who gifted the subscription. Null if anonymous.
	GifterUserID string `json:"gifter_user_id,omitempty"`

	// Optional. The user name of the user who gifted the subscription. Null if anonymous.
	GifterUserName string `json:"gifter_user_name,omitempty"`

	// Optional. The user login of the user who gifted the subscription. Null if anonymous.
	GifterUserLogin string `json:"gifter_user_login,omitempty"`
}

// ChatNotificationAnnouncementEvent represents information about the announcement event.
type ChatNotificationAnnouncementEvent struct {
	// Color of the announcement.
	Color string `json:"color"`
}

// ChatNotificationCharityDonationEvent represents information about the charity donation event.
type ChatNotificationCharityDonationEvent struct {
	// Name of the charity.
	CharityName string `json:"charity_name"`

	// An object that contains the amount of money that the user paid.
	Amount ChatNotificationCharityDonationEventDonationAmount `json:"amount"`
}

// ChatNotificationCharityDonationEventDonationAmount represents the amount of money that the user paid.
type ChatNotificationCharityDonationEventDonationAmount struct {
	// The monetary amount. The amount is specified in the currency’s minor unit.
	// For example, the minor units for USD is cents, so if the amount is $5.50 USD, value is set to 550.
	Value int `json:"value"`

	// The number of decimal places used by the currency.
	// For example, USD uses two decimal places.
	DecimalPlaces int `json:"decimal_places"`

	// The ISO-4217 three-letter currency code that identifies the type of currency in value.
	Currency string `json:"currency"`
}

// ChatNotificationBitsBadgeTierEvent represents information about the bits badge tier event.
type ChatNotificationBitsBadgeTierEvent struct {
	// The tier of the Bits badge the user just earned. For example, 100, 1000, or 10000.
	Tier int `json:"tier"`
}

// ChatReply represents metadata about a reply message.
type ChatReply struct {
	// An ID that uniquely identifies the parent message that this message is replying to.
	ParentMessageID string `json:"parent_message_id"`
	// The message body of the parent message.
	ParentMessageBody string `json:"parent_message_body"`
	// User ID of the sender of the parent message.
	ParentUserID string `json:"parent_user_id"`
	// User name of the sender of the parent message.
	ParentUserName string `json:"parent_user_name"`
	// User login of the sender of the parent message.
	ParentUserLogin string `json:"parent_user_login"`
	// An ID that identifies the parent message of the reply thread.
	ThreadMessageID string `json:"thread_message_id"`
	// User ID of the sender of the thread’s parent message.
	ThreadUserID string `json:"thread_user_id"`
	// User name of the sender of the thread’s parent message.
	ThreadUserName string `json:"thread_user_name"`
	// User login of the sender of the thread’s parent message.
	ThreadUserLogin string `json:"thread_user_login"`
}

// ChatCheer represents metadata about a message cheer.
type ChatCheer struct {
	// The amount of Bits the user cheered.
	Bits int `json:"bits"`
}

package beta

type EventChannelUnbanRequestCreate struct {
	// The ID of the unban request.
	ID string `json:"id"`
	// The broadcaster’s user ID for the channel the unban request was created for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster’s login name.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// User ID of user that is requesting to be unbanned.
	UserID string `json:"user_id"`
	// The user’s login name.
	UserLogin string `json:"user_login"`
	// The user’s display name.
	UserName string `json:"user_name"`
	// Message sent in the unban request.
	Text string `json:"text"`
	// The UTC timestamp (in RFC3339 format) of when the unban request was created.
	CreatedAt string `json:"created_at"`
}

type EventChannelUnbanRequestResolve struct {
	// The ID of the unban request.
	ID string `json:"id"`
	// The broadcaster’s user ID for the channel the unban request was created for.
	BroadcasterUserID string `json:"broadcaster_user_id"`
	// The broadcaster’s login name.
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	// The broadcaster’s display name.
	BroadcasterUserName string `json:"broadcaster_user_name"`
	// Optional. User ID of moderator who approved/denied the request.
	ModeratorUserID string `json:"moderator_user_id"`
	// Optional. The moderator’s login name
	ModeratorUserLogin string `json:"moderator_user_login"`
	// Optional. The moderator’s display name
	ModeratorUserName string `json:"moderator_user_name"`
	// User ID of user that requested to be unbanned.
	UserID string `json:"user_id"`
	// 	The user’s login name.
	UserLogin string `json:"user_login"`
	// The user’s display name.
	UserName string `json:"user_name"`
	// Optional. Resolution text supplied by the mod/broadcaster upon approval/denial of the request.
	ResolutionText string `json:"resolution_text"`
	// Dictates whether the unban request was approved or denied.
	// Can be the following: approved, canceled, denied
	Status ChannelUnbanRequestResolveStatus `json:"status"`
}

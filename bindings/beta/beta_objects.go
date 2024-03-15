package beta

type ChannelUnbanRequestResolveStatus string

const (
	ChannelUnbanRequestResolveStatusApproved ChannelUnbanRequestResolveStatus = "approved"
	ChannelUnbanRequestResolveStatusCanceled ChannelUnbanRequestResolveStatus = "canceled"
	ChannelUnbanRequestResolveStatusDenied   ChannelUnbanRequestResolveStatus = "denied"
)

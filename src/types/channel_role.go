package types

type ChannelRole struct {
	ID              string `json:"id"`
	ChannelId       string `json:"channel_id"`
	AppserverId     string `json:"appserver_id"`
	AppserverRoleId string `json:"appserver_role_id"`
}

type ChannelRoleCreate struct {
	ChannelId       string `json:"channel_id"`
	AppserverId     string `json:"appserver_id"`
	AppserverRoleId string `json:"appserver_role_id"`
}

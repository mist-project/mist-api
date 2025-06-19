package types

type Channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

type ChannelCreate struct {
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
	IsPrivate   bool   `json:"is_private,omitempty"`
}

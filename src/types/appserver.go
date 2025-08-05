package types

type Appserver struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"is_owner"`
}

type AppserverDetail struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	IsOwner  bool            `json:"is_owner"`
	Roles    []AppserverRole `json:"roles"`
	Channels []Channel       `json:"channels"`
}

type AppserverCreate struct {
	Name string `json:"name"`
}

type AppserverAndSub struct {
	Appserver Appserver `json:"appserver"`
	SubId     string    `json:"sub_id"`
}

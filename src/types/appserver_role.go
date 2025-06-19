package types

type AppserverRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

type AppserverRoleCreate struct {
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

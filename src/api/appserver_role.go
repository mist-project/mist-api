package api

type AppserverRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

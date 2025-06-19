package types

type AppserverSub struct {
	ID          string `json:"id"`
	AppuserId   string `json:"appuser_id"`
	AppserverId string `json:"appserver_id"`
}

type AppserverSubCreate struct {
	AppserverId string `json:"appserver_id"`
}

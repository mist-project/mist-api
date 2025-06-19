package types

type Appuser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AppuserAppserverSub struct {
	Appuser Appuser `json:"appuser"`
	SubId   string  `json:"appserver_sub_id"`
}

package types

type AppserverRoleSub struct {
	ID              string `json:"id"`
	AppuserId       string `json:"appuser_id"`
	AppserverRoleId string `json:"appserver_role_id"`
	AppserverId     string `json:"appserver_id"`
}

type AppserverRoleSubCreate struct {
	AppuserId       string `json:"appuser_id"`
	AppserverRoleId string `json:"appserver_role_id"`
	AppserverId     string `json:"appserver_id"`
	AppserverSubId  string `json:"appserver_sub_id"`
}

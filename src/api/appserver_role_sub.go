package api

import (
	"net/http"

	"mistapi/src/auth"
	pb_appserver_role_sub "mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

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

func appserverRoleSubRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", AppserverRoleSubListHandler)          // list role subs for server
	r.Post("/", AppserverRoleSubCreateHandler)       // create a new role sub
	r.Delete("/{id}", AppserverRoleSubDeleteHandler) // delete a role sub
	return r
}

// AppserverRoleSubCreateHandler godoc
// @Summary      Create a role subscription for a user
// @Description  Assign a role to a user in a server subscription
// @Tags         appserver-role-subs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver_role_sub  body  AppserverRoleSubCreate  true  "AppserverRoleSubCreate"
// @Success      204
// @Router       /api/v1/appserver-role-subs [post]
func AppserverRoleSubCreateHandler(w http.ResponseWriter, r *http.Request) {
	var roleSub AppserverRoleSubCreate

	err := DecodeRequestBody(w, r, &roleSub)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err = c.GetAppserverRoleSubClient().Create(
		ctx, &pb_appserver_role_sub.CreateRequest{
			AppuserId:       roleSub.AppuserId,
			AppserverRoleId: roleSub.AppserverRoleId,
			AppserverId:     roleSub.AppserverId,
			AppserverSubId:  roleSub.AppserverSubId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.NoContent(w, r)
}

// AppserverRoleSubListHandler godoc
// @Summary      List all role subscriptions in a server
// @Description  Get all user role subscriptions in a given server
// @Tags         appserver-role-subs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver_id  query  string  true  "Appserver ID"
// @Success      200  {array}  AppserverRoleSub
// @Router       /api/v1/appserver-role-subs [get]
func AppserverRoleSubListHandler(w http.ResponseWriter, r *http.Request) {
	appserverID := r.URL.Query().Get("appserver_id")

	if appserverID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, CreateErrorResponse("Appserver ID is required"))
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverRoleSubClient().ListServerRoleSubs(
		ctx, &pb_appserver_role_sub.ListServerRoleSubsRequest{
			AppserverId: appserverID,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	res := make([]AppserverRoleSub, 0, len(response.AppserverRoleSubs))

	for _, a := range response.AppserverRoleSubs {
		res = append(res, AppserverRoleSub{
			ID:              a.Id,
			AppuserId:       a.AppuserId,
			AppserverRoleId: a.AppserverRoleId,
			AppserverId:     a.AppserverId,
		})
	}

	render.JSON(w, r, CreateResponse(res))
}

// AppserverRoleSubDeleteHandler godoc
// @Summary      Delete a user role subscription
// @Description  Delete a role sub entry by ID
// @Tags         appserver-role-subs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Role Sub ID"
// @Success      204
// @Router       /api/v1/appserver-role-subs/{id} [delete]
func AppserverRoleSubDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetAppserverRoleSubClient().Delete(
		ctx, &pb_appserver_role_sub.DeleteRequest{
			Id: id,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

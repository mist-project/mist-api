package api

import (
	"net/http"

	"mistapi/src/auth"
	"mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/service"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func appserverRoleSubRouter() http.Handler {
	r := chi.NewRouter()

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
// @Param        appserver_role_sub  body  types.AppserverRoleSubCreate  true  "AppserverRoleSubCreate"
// @Success      204
// @Router       /api/v1/appserver-role-subs [post]
func AppserverRoleSubCreateHandler(w http.ResponseWriter, r *http.Request) {
	var roleSub types.AppserverRoleSubCreate

	err := DecodeRequestBody(w, r, &roleSub)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err = c.GetAppserverRoleSubClient().Create(
		ctx, &appserver_role_sub.CreateRequest{
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
		ctx, &appserver_role_sub.DeleteRequest{
			Id: id,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

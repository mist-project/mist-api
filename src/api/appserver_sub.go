package api

import (
	"net/http"

	"mistapi/src/auth"
	"mistapi/src/protos/v1/appserver_sub"
	"mistapi/src/service"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func appserverSubRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", AppserverSubCreateHandler)       // create an appserver sub
	r.Delete("/{id}", AppserverSubDeleteHandler) // delete an appserver sub
	return r
}

// AppserverSubCreateHandler godoc
// @Summary      Create an appserver sub
// @Description  Create an appserver sub
// @Tags         appserver-subs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver  body      types.AppserverSubCreate  true  "AppserverSubCreate"
// @Success      201 {object} types.AppserverSub
// @Router       /api/v1/appserver-subs [post]
func AppserverSubCreateHandler(w http.ResponseWriter, r *http.Request) {
	var sub types.AppserverSubCreate

	err := DecodeRequestBody(w, r, &sub)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverSubClient().Create(
		ctx, &appserver_sub.CreateRequest{
			AppserverId: sub.AppserverId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse(&types.AppserverSub{
		ID:          response.AppserverSub.Id,
		AppserverId: response.AppserverSub.AppserverId,
	}))
}

// AppserverSubDeleteHandler godoc
// @Summary      Delete appserver sub by id
// @Description  Delete appserver sub by id (removing a user from channel)
// @Tags         appserver-subs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver sub ID"
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/appserver-subs/{id} [delete]
func AppserverSubDeleteHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetAppserverSubClient().Delete(
		ctx, &appserver_sub.DeleteRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

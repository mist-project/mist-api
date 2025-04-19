package api

import (
	"net/http"

	"mistapi/src/auth"
	pb "mistapi/src/protos/v1/gen"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AppserverSub struct {
	ID          string `json:"id"`
	AppuserId   string `json:"appuser_id"`
	AppserverId string `json:"appserver_id"`
}

type AppserverSubCreate struct {
	AppserverId string `json:"appserver_id"`
}

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
// @Param        appserver  body      AppserverSubCreate  true  "AppserverSubCreate"
// @Success      201 {object} AppserverSub
// @Router       /api/v1/appserver-subs [post]
func AppserverSubCreateHandler(w http.ResponseWriter, r *http.Request) {
	var sub AppserverSubCreate

	err := DecodeRequestBody(w, r, &sub)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().CreateAppserverSub(
		ctx, &pb.CreateAppserverSubRequest{
			AppserverId: sub.AppserverId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse(&AppserverSub{
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
	_, err := c.GetServerClient().DeleteAppserverSub(
		ctx, &pb.DeleteAppserverSubRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

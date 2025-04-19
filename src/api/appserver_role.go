package api

import (
	"net/http"

	"mistapi/src/auth"
	pb "mistapi/src/protos/v1/gen"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AppserverRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

type AppserverRoleCreate struct {
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

func appserverRoleRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", AppserverRoleCreateHandler)       // create an appserver role
	r.Delete("/{id}", AppserverRoleDeleteHandler) // delete an appserver role
	return r
}

// AppserverRoleCreateHandler godoc
// @Summary      Create an appserver role
// @Description  Create an appserver role
// @Tags         appserver-roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver  body      AppserverRoleCreate  true  "AppserverRoleCreate"
// @Success      201 {object} AppserverRole
// @Router       /api/v1/appserver-roles [post]
func AppserverRoleCreateHandler(w http.ResponseWriter, r *http.Request) {
	var role AppserverRoleCreate

	err := DecodeRequestBody(w, r, &role)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().CreateAppserverRole(
		ctx, &pb.CreateAppserverRoleRequest{
			Name:        role.Name,
			AppserverId: role.AppserverId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse(&AppserverRole{
		ID:          response.AppserverRole.Id,
		Name:        response.AppserverRole.Name,
		AppserverId: response.AppserverRole.AppserverId,
	}))
}

// AppserverRoleDeleteHandler godoc
// @Summary      Delete appserver role by id
// @Description  Delete appserver role by id, only owners of server can perform this action
// @Tags         appserver-roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver role ID"
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/appserver-roles/{id} [delete]
func AppserverRoleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetServerClient().DeleteAppserverRole(
		ctx, &pb.DeleteAppserverRoleRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

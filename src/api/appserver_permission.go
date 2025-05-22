package api

import (
	"net/http"

	"mistapi/src/auth"
	pb_appserver_permission "mistapi/src/protos/v1/appserver_permission"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AppserverPermission struct {
	ID          string `json:"id"`
	AppuserId   string `json:"appuser_id"`
	AppserverId string `json:"appserver_id"`
}

type AppserverPermissionCreate struct {
	AppuserId   string `json:"appuser_id"`
	AppserverId string `json:"appserver_id"`
}

func appserverPermissionRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", AppserverPermissionListHandler)          // get all permission roles for server
	r.Post("/", AppserverPermissionCreateHandler)       // create an appserver permission
	r.Delete("/{id}", AppserverPermissionDeleteHandler) // delete an appserver permission
	return r
}

// AppserverPermissionCreateHandler godoc
// @Summary      Create an permission role for user in server
// @Description  Create an permission role for user in server
// @Tags         appserver-permissions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver  body      AppserverPermissionCreate  true  "AppserverPermissionCreate"
// @Success      204
// @Router       /api/v1/appserver-permissions [post]
func AppserverPermissionCreateHandler(w http.ResponseWriter, r *http.Request) {
	var role AppserverPermissionCreate

	err := DecodeRequestBody(w, r, &role)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err = c.GetAppserverPermissionClient().Create(
		ctx, &pb_appserver_permission.CreateRequest{
			AppuserId:   role.AppuserId,
			AppserverId: role.AppserverId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.NoContent(w, r)
}

// List godoc
// @Summary      List of all users with permission role for server
// @Description  List of all users with permission role for server
// @Tags         appserver-permissions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver_id  query      string  true  "Appserver ID"
// @Success      200  {array}  Appserver
// @Router       /api/v1/appserver-permissions [get]
func AppserverPermissionListHandler(w http.ResponseWriter, r *http.Request) {
	appserverID := r.URL.Query().Get("appserver_id")

	if appserverID == "" {
		// If appserverid is missing, return a bad request error
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, CreateErrorResponse("Appserver ID is required"))
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverPermissionClient().ListAppserverUsers(
		ctx, &pb_appserver_permission.ListAppserverUsersRequest{
			AppserverId: appserverID,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	res := make([]AppserverPermission, 0, len(response.AppserverPermissions))

	for _, a := range response.AppserverPermissions {
		res = append(res, AppserverPermission{
			ID:          a.Id,
			AppuserId:   a.AppuserId,
			AppserverId: a.AppserverId,
		})
	}

	render.JSON(w, r, CreateResponse(res))
}

// AppserverPermissionDeleteHandler godoc
// @Summary      Delete appserver role by id
// @Description  Delete appserver role by id, only owners of server can perform this action
// @Tags         appserver-permissions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver role ID"
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/appserver-permissions/{id} [delete]
func AppserverPermissionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetAppserverPermissionClient().Delete(
		ctx, &pb_appserver_permission.DeleteRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

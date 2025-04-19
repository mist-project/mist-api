package api

import (
	"fmt"
	"net/http"

	"mistapi/src/auth"
	pb "mistapi/src/protos/v1/gen"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func appserverRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", AppserverCreateHandler) // create an appserver
	r.Get("/", AppserverListHandler)    // list all existing servers (most likely to be deprecated)

	r.Get("/{id}", AppserverDetailHandler)          // get all appserver details
	r.Get("/{id}/subs", AppserverListSubsHandler)   // get all user subscriptions appserver has
	r.Get("/{id}/roles", AppserverListRolesHandler) // get all roles an appserver has
	r.Delete("/{id}", AppserverDeleteHandler)       // delete an appserver

	r.Get("/user-subs", AppserverUserListSubsHandler) // get all appservers a user is subscribed to
	return r
}

type Appserver struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"is_owner"`
}

type AppserverDetail struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"is_owner"`
}

type AppserverCreate struct {
	Name string `json:"name"`
}

type AppserverAndSub struct {
	Appserver Appserver `json:"appserver"`
	SubId     string    `json:"sub_id"`
}

// AppserverCreateHandler godoc
// @Summary      Create an appserver
// @Description  Create an appserver
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver  body      AppserverCreate  true  "AppserverCreate"
// @Success      201 {object} Appserver
// @Router       /api/v1/appserver [post]
func AppserverCreateHandler(w http.ResponseWriter, r *http.Request) {
	var appserver AppserverCreate

	err := DecodeRequestBody(w, r, &appserver)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().CreateAppserver(
		ctx, &pb.CreateAppserverRequest{
			Name: appserver.Name,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse(response.Appserver))
}

// List godoc
// @Summary      List of all appservers
// @Description  List of all appservers
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  Appserver
// @Router       /api/v1/appserver [get]
func AppserverListHandler(w http.ResponseWriter, r *http.Request) {
	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().ListAppservers(
		ctx, &pb.ListAppserversRequest{},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	res := make([]Appserver, 0, len(response.Appservers))

	for _, a := range response.Appservers {
		res = append(res, Appserver{
			ID:      a.Id,
			Name:    a.Name,
			IsOwner: a.IsOwner,
		})
	}

	render.JSON(w, r, CreateResponse(res))
}

// AppserverUserListSubsHandler godoc
// @Summary      List of appservers for a particular user
// @Description  List of appservers for a particular user (user in jwt token)
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  AppserverAndSub
// @Router       /api/v1/appserver/subs [get]
func AppserverUserListSubsHandler(w http.ResponseWriter, r *http.Request) {
	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().GetUserAppserverSubs(
		ctx, &pb.GetUserAppserverSubsRequest{},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	res := make([]AppserverAndSub, 0, len(response.Appservers))

	for _, a := range response.Appservers {
		res = append(res, AppserverAndSub{
			Appserver: Appserver{
				ID:      a.Appserver.Id,
				Name:    a.Appserver.Name,
				IsOwner: a.Appserver.IsOwner,
			},
			SubId: a.SubId,
		})
	}

	render.JSON(w, r, CreateResponse(res))
}

// AppserverDetailHandler godoc
// @Summary      Gets all details of an appserver
// @Description  Gets (almost) everything related to an appserver, except its user subscriptions
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver ID"
// @Security     BearerAuth
// @Success      200 {array} AppserverDetail
// @Router       /api/v1/appserver/{id} [get]
func AppserverDetailHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")
	fmt.Printf("->> %s|\n", sId)

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().GetByIdAppserver(
		ctx, &pb.GetByIdAppserverRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.JSON(w, r, CreateResponse(&AppserverDetail{
		ID:      response.Appserver.Id,
		Name:    response.Appserver.Name,
		IsOwner: response.Appserver.IsOwner,
	}))
}

// AppserverListSubsHandler godoc
// @Summary      Gets all user subs to a server
// @Description  Gets all users in the server and their sub id
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver ID"
// @Security     BearerAuth
// @Success      200 {array} AppuserAppserverSub
// @Router       /api/v1/appserver/{id}/subs [get]
func AppserverListSubsHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().GetAllUsersAppserverSubs(
		ctx, &pb.GetAllUsersAppserverSubsRequest{
			AppserverId: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	subs := make([]AppuserAppserverSub, 0, len(response.Appusers))

	for _, sub := range response.Appusers {
		subs = append(subs, AppuserAppserverSub{
			Appuser: Appuser{ID: sub.Appuser.Id, Username: sub.Appuser.Username},
			SubId:   sub.SubId,
		})
	}

	render.JSON(w, r, CreateResponse(subs))
}

// AppserverListSubsHandler godoc
// @Summary      Gets all roles in a appserver
// @Description  Gets all roles in an appserver
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver ID"
// @Security     BearerAuth
// @Success      200 {array} AppserverRole
// @Router       /api/v1/appserver/{id}/roles [get]
func AppserverListRolesHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetServerClient().GetAllAppserverRoles(
		ctx, &pb.GetAllAppserverRolesRequest{
			AppserverId: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	roles := make([]AppserverRole, 0, len(response.AppserverRoles))

	for _, role := range response.AppserverRoles {
		roles = append(roles, AppserverRole{
			ID:          role.Id,
			Name:        role.Name,
			AppserverId: role.AppserverId,
		})
	}

	render.JSON(w, r, CreateResponse(roles))
}

// AppserverDeleteHandler godoc
// @Summary      Delete Appserver by id
// @Description  Delete an appserver, only owners of server can perform this action
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver ID"
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/appserver/{id} [delete]
func AppserverDeleteHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetServerClient().DeleteAppserver(
		ctx, &pb.DeleteAppserverRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

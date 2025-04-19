package api

import (
	"net/http"

	"mistapi/src/auth"
	pb "mistapi/src/protos/v1/gen"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// A completely separate router for administrator routes
func appserverRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.AuthenticateMiddleware)
	r.Get("/", List)
	r.Get("/subs", ListSubs)
	// r.Get("/{id}", getUser)
	r.Post("/", Create)
	r.Delete("/{id}", Delete)
	return r
}

type Appserver struct {
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

// Create godoc
// @Summary      Create Appserver
// @Description  Create an appserver
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver  body      AppserverCreate  true  "AppserverCreate"
// @Success      201 {object} Appserver
// @Router       /api/v1/appserver [post]
func Create(w http.ResponseWriter, r *http.Request) {
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
// @Summary      List Appservers
// @Description  Get a list of appservers
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  Appserver
// @Router       /api/v1/appserver [get]
func List(w http.ResponseWriter, r *http.Request) {
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

// list/subs godoc
// @Summary      List Appserver and SubId
// @Description  Get a list of appservers for a particular user (user in jwt token0)
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  AppserverAndSub
// @Router       /api/v1/appserver/subs [get]
func ListSubs(w http.ResponseWriter, r *http.Request) {
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

// Delete godoc
// @Summary      Delete Appserver by id
// @Description  Delete an appserver, only owners of server can perform this action
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Server ID"
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/appserver/{id} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
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

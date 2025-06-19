package api

import (
	"net/http"

	"mistapi/src/auth"
	"mistapi/src/protos/v1/appserver"
	"mistapi/src/protos/v1/appserver_role"
	"mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/protos/v1/appserver_sub"
	"mistapi/src/protos/v1/channel"
	"mistapi/src/protos/v1/channel_role"
	"mistapi/src/service"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func appserverRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", AppserverCreateHandler) // create an appserver
	r.Get("/", AppserverListHandler)    // list all existing servers (most likely to be deprecated)

	r.Get("/{id}", AppserverDetailHandler)                                     // get all appserver details
	r.Get("/{id}/channels", AppserverListChannelsHandler)                      // get all channels in a server
	r.Get("/{sid}/channels/{cid}/channel-roles", AppserverChannelRolesHandler) // get all channel roles in a server
	r.Get("/{id}/subs", AppserverListSubsHandler)                              // get all appserver user subscriptions
	r.Get("/{id}/roles", AppserverListRolesHandler)                            // get all appserver roles
	r.Get("/{id}/role-subs", AppserverListRoleSubHandler)                      // get all appservers' role subscriptions

	r.Delete("/{id}", AppserverDeleteHandler) // delete an appserver

	return r
}

// AppserverCreateHandler godoc
// @Summary      Create an appserver
// @Description  Create an appserver
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver  body      types.AppserverCreate  true  "AppserverCreate"
// @Success      201 {object} types.Appserver
// @Router       /api/v1/appservers [post]
func AppserverCreateHandler(w http.ResponseWriter, r *http.Request) {
	var s types.AppserverCreate

	err := DecodeRequestBody(w, r, &s)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverClient().Create(
		ctx, &appserver.CreateRequest{
			Name: s.Name,
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
// @Summary      List of all appservers for a particular user
// @Description  List of all appservers for a particular user (user in jwt token)
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  types.Appserver
// @Router       /api/v1/appservers [get]
func AppserverListHandler(w http.ResponseWriter, r *http.Request) {
	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverSubClient().ListUserServerSubs(
		ctx, &appserver_sub.ListUserServerSubsRequest{},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	res := make([]types.AppserverAndSub, 0, len(response.Appservers))

	for _, a := range response.Appservers {
		res = append(res, types.AppserverAndSub{
			Appserver: types.Appserver{
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
// @Success      200 {array} types.AppserverDetail
// @Router       /api/v1/appservers/{id} [get]
func AppserverDetailHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverClient().GetById(
		ctx, &appserver.GetByIdRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.JSON(w, r, CreateResponse(&types.AppserverDetail{
		ID:      response.Appserver.Id,
		Name:    response.Appserver.Name,
		IsOwner: response.Appserver.IsOwner,
	}))
}

// AppserverListSubsHandler godoc
// @Summary      Gets all user subscribed to a server
// @Description  Gets all users in the server and their sub id
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Appserver ID"
// @Security     BearerAuth
// @Success      200 {array} types.AppuserAppserverSub
// @Router       /api/v1/appservers/{id}/subs [get]
func AppserverListSubsHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverSubClient().ListAppserverUserSubs(
		ctx, &appserver_sub.ListAppserverUserSubsRequest{
			AppserverId: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	subs := make([]types.AppuserAppserverSub, 0, len(response.Appusers))

	for _, sub := range response.Appusers {
		subs = append(subs, types.AppuserAppserverSub{
			Appuser: types.Appuser{ID: sub.Appuser.Id, Username: sub.Appuser.Username},
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
// @Success      200 {array} types.AppserverRole
// @Router       /api/v1/appservers/{id}/roles [get]
func AppserverListRolesHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverRoleClient().ListServerRoles(
		ctx, &appserver_role.ListServerRolesRequest{
			AppserverId: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	roles := make([]types.AppserverRole, 0, len(response.AppserverRoles))

	for _, role := range response.AppserverRoles {
		roles = append(roles, types.AppserverRole{
			ID:          role.Id,
			Name:        role.Name,
			AppserverId: role.AppserverId,
		})
	}

	render.JSON(w, r, CreateResponse(roles))
}

// AppserverListRoleSubHandler godoc
// @Summary      List all role subscriptions in a server
// @Description  Get all user role subscriptions in a given server
// @Tags         appserver-role-subs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Appserver ID"
// @Success      200  {array}  types.AppserverRoleSub
// @Router       /api/v1/appservers/{id}/appserver-role-subs [get]
func AppserverListRoleSubHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetAppserverRoleSubClient().ListServerRoleSubs(
		ctx, &appserver_role_sub.ListServerRoleSubsRequest{
			AppserverId: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	res := make([]types.AppserverRoleSub, 0, len(response.AppserverRoleSubs))

	for _, a := range response.AppserverRoleSubs {
		res = append(res, types.AppserverRoleSub{
			ID:              a.Id,
			AppuserId:       a.AppuserId,
			AppserverRoleId: a.AppserverRoleId,
			AppserverId:     a.AppserverId,
		})
	}

	render.JSON(w, r, CreateResponse(res))
}

// AppserverListChannelsHandler godoc
// @Summary      List channels for a given appserver ID
// @Description  List all channels associated with a specific appserver ID
// @Tags         channel
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Appserver ID"
// @Success      200          {array}   types.AppserverSub
// @Failure      400          {object}  ErrorResponse "Invalid appserver ID"
// @Failure      500          {object}  ErrorResponse "Internal Server Error"
// @Router       /api/v1/appservers/{id}/channels [get]
func AppserverListChannelsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the appserver ID from URL parameters

	sId := chi.URLParam(r, "id")

	// Authorization and gRPC context setup
	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	// Create a new gRPC client and make the request to list channels for the appserver
	c := service.NewGrpcClient()
	response, err := c.GetChannelClient().ListServerChannels(
		ctx, &channel.ListServerChannelsRequest{
			AppserverId: sId,
		},
	)

	if err != nil {
		// Handle gRPC error and return it as a response
		HandleGrpcError(w, r, err)
		return
	}

	channels := make([]types.Channel, 0, len(response.Channels))

	for _, c := range response.Channels {
		channels = append(channels, types.Channel{
			ID:          c.Id,
			Name:        c.Name,
			AppserverId: c.AppserverId,
		})
	}
	// Successfully fetched channels, return them in the response
	render.JSON(w, r, CreateResponse(channels))
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
// @Router       /api/v1/appservers/{id} [delete]
func AppserverDeleteHandler(w http.ResponseWriter, r *http.Request) {
	sId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetAppserverClient().Delete(
		ctx, &appserver.DeleteRequest{
			Id: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

// AppserverChannelRolesHandler godoc
// @Summary      List all roles assigned to a channel
// @Description  Get all server roles mapped to a specific channel
// @Tags         channel
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        channel_id    path  string  true  "Channel ID"
// @Param        appserver_id  path  string  true  "Appserver ID"
// @Success      200  {array}  types.ChannelRole
// @Router       /api/v1/appservers/{sid}/channels/{cid}/channel-roles [get]
func AppserverChannelRolesHandler(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "cid")
	sId := chi.URLParam(r, "sid")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	res, err := c.GetChannelRoleClient().ListChannelRoles(
		ctx, &channel_role.ListChannelRolesRequest{
			ChannelId:   channelID,
			AppserverId: sId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	response := make([]types.ChannelRole, 0, len(res.ChannelRoles))
	for _, r := range res.ChannelRoles {
		response = append(response, types.ChannelRole{
			ID:              r.Id,
			ChannelId:       r.ChannelId,
			AppserverId:     r.AppserverId,
			AppserverRoleId: r.AppserverRoleId,
		})
	}

	render.JSON(w, r, CreateResponse(response))
}

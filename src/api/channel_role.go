package api

import (
	"net/http"

	"mistapi/src/auth"
	pb_channel_role "mistapi/src/protos/v1/channel_role"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ChannelRole struct {
	ID              string `json:"id"`
	ChannelId       string `json:"channel_id"`
	AppserverId     string `json:"appserver_id"`
	AppserverRoleId string `json:"appserver_role_id"`
}

type ChannelRoleCreate struct {
	ChannelId       string `json:"channel_id"`
	AppserverId     string `json:"appserver_id"`
	AppserverRoleId string `json:"appserver_role_id"`
}

func channelRoleRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", ChannelRoleListHandler)          // list channel roles
	r.Post("/", ChannelRoleCreateHandler)       // create a channel role
	r.Delete("/{id}", ChannelRoleDeleteHandler) // delete a channel role
	return r
}

// ChannelRoleCreateHandler godoc
// @Summary      Create a role for a channel
// @Description  Assign a server role to a channel
// @Tags         channel-roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        channel_role  body  ChannelRoleCreate  true  "ChannelRoleCreate"
// @Success      204
// @Router       /api/v1/channel-roles [post]
func ChannelRoleCreateHandler(w http.ResponseWriter, r *http.Request) {
	var role ChannelRoleCreate

	err := DecodeRequestBody(w, r, &role)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err = c.GetChannelRoleClient().Create(
		ctx, &pb_channel_role.CreateRequest{
			ChannelId:       role.ChannelId,
			AppserverId:     role.AppserverId,
			AppserverRoleId: role.AppserverRoleId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.NoContent(w, r)
}

// ChannelRoleListHandler godoc
// @Summary      List all roles assigned to a channel
// @Description  Get all server roles mapped to a specific channel
// @Tags         channel-roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        channel_id    query  string  true  "Channel ID"
// @Param        appserver_id  query  string  true  "Appserver ID"
// @Success      200  {array}  ChannelRole
// @Router       /api/v1/channel-roles [get]
func ChannelRoleListHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel_id")
	appserverID := r.URL.Query().Get("appserver_id")

	if channelID == "" || appserverID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, CreateErrorResponse("Channel ID and Appserver ID are required"))
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	res, err := c.GetChannelRoleClient().ListChannelRoles(
		ctx, &pb_channel_role.ListChannelRolesRequest{
			ChannelId:   channelID,
			AppserverId: appserverID,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	response := make([]ChannelRole, 0, len(res.ChannelRoles))
	for _, r := range res.ChannelRoles {
		response = append(response, ChannelRole{
			ID:              r.Id,
			ChannelId:       r.ChannelId,
			AppserverId:     r.AppserverId,
			AppserverRoleId: r.AppserverRoleId,
		})
	}

	render.JSON(w, r, CreateResponse(response))
}

// ChannelRoleDeleteHandler godoc
// @Summary      Delete a channel role
// @Description  Delete a role assigned to a channel
// @Tags         channel-roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Channel Role ID"
// @Success      204
// @Router       /api/v1/channel-roles/{id} [delete]
func ChannelRoleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetChannelRoleClient().Delete(ctx, &pb_channel_role.DeleteRequest{Id: id})

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

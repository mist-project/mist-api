package api

import (
	"net/http"

	"mistapi/src/auth"
	"mistapi/src/protos/v1/channel_role"
	"mistapi/src/service"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func channelRoleRouter() http.Handler {
	r := chi.NewRouter()

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
// @Param        channel_role  body  types.ChannelRoleCreate  true  "ChannelRoleCreate"
// @Success      204
// @Router       /api/v1/channel-roles [post]
func ChannelRoleCreateHandler(w http.ResponseWriter, r *http.Request) {
	var role types.ChannelRoleCreate

	err := DecodeRequestBody(w, r, &role)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err = c.GetChannelRoleClient().Create(
		ctx, &channel_role.CreateRequest{
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
	_, err := c.GetChannelRoleClient().Delete(ctx, &channel_role.DeleteRequest{Id: id})

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

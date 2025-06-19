package api

import (
	"net/http"

	"mistapi/src/auth"
	"mistapi/src/protos/v1/channel"
	"mistapi/src/service"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func channelRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", ChannelCreateHandler)       // create a channel
	r.Delete("/{id}", ChannelDeleteHandler) // delete a channel
	return r
}

// ChannelCreateHandler godoc
// @Summary      Create a channel in a server
// @Description  Create a channel in a server
// @Tags         channel
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        channel  body      types.ChannelCreate  true  "ChannelCreate"
// @Success      201 {object} types.Channel
// @Router       /api/v1/channels [post]
func ChannelCreateHandler(w http.ResponseWriter, r *http.Request) {
	var c types.ChannelCreate

	err := DecodeRequestBody(w, r, &c)

	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	client := service.NewGrpcClient()
	response, err := client.GetChannelClient().Create(
		ctx, &channel.CreateRequest{
			Name:        c.Name,
			AppserverId: c.AppserverId,
			IsPrivate:   c.IsPrivate,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse(&types.Channel{
		ID:          response.Channel.Id,
		Name:        response.Channel.Name,
		AppserverId: response.Channel.AppserverId,
	}))
}

// ChannelDeleteHandler godoc
// @Summary      Delete channel by id
// @Description  Delete channel by id
// @Tags         channel
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Channel ID"
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/channels/{id} [delete]
func ChannelDeleteHandler(w http.ResponseWriter, r *http.Request) {
	cId := chi.URLParam(r, "id")

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	_, err := c.GetChannelClient().Delete(
		ctx, &channel.DeleteRequest{
			Id: cId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

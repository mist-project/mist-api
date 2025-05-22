package api

import (
	"net/http"

	"mistapi/src/auth"
	pb_channel "mistapi/src/protos/v1/channel"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

type ChannelCreate struct {
	Name        string `json:"name"`
	AppserverId string `json:"appserver_id"`
}

func channelRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", ChannelCreateHandler)       // create a channel
	r.Get("/", ListChannelsHandler)         // get all channels in a server
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
// @Param        channel  body      ChannelCreate  true  "ChannelCreate"
// @Success      201 {object} Channel
// @Router       /api/v1/channels [post]
func ChannelCreateHandler(w http.ResponseWriter, r *http.Request) {
	var channel ChannelCreate

	err := DecodeRequestBody(w, r, &channel)
	if err != nil {
		return
	}

	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	c := service.NewGrpcClient()
	response, err := c.GetChannelClient().Create(
		ctx, &pb_channel.CreateRequest{
			Name:        channel.Name,
			AppserverId: channel.AppserverId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse(&Channel{
		ID:          response.Channel.Id,
		Name:        response.Channel.Name,
		AppserverId: response.Channel.AppserverId,
	}))
}

// ListChannelsHandler godoc
// @Summary      List channels for a given appserver ID
// @Description  List all channels associated with a specific appserver ID
// @Tags         channel
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appserver_id  query      string  true  "Appserver ID"
// @Success      200          {array}   AppserverSub
// @Failure      400          {object}  ErrorResponse "Invalid appserver ID"
// @Failure      500          {object}  ErrorResponse "Internal Server Error"
// @Router       /api/v1/channels [get]
func ListChannelsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the appserver ID from URL parameters

	appserverID := r.URL.Query().Get("appserver_id")

	if appserverID == "" {
		// If appserverid is missing, return a bad request error
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, CreateErrorResponse("Appserver ID is required"))
		return
	}

	// Authorization and gRPC context setup
	authT, _ := auth.GetAuthotizationToken(r)
	ctx, cancel := service.SetupGrpcHeaders(authT.Token)
	defer cancel()

	// Create a new gRPC client and make the request to list channels for the appserver
	c := service.NewGrpcClient()
	response, err := c.GetChannelClient().ListServerChannels(
		ctx, &pb_channel.ListServerChannelsRequest{
			AppserverId: appserverID,
		},
	)

	if err != nil {
		// Handle gRPC error and return it as a response
		HandleGrpcError(w, r, err)
		return
	}

	channels := make([]Channel, 0, len(response.Channels))

	for _, c := range response.Channels {
		channels = append(channels, Channel{
			ID:          c.Id,
			Name:        c.Name,
			AppserverId: c.AppserverId,
		})
	}
	// Successfully fetched channels, return them in the response
	render.JSON(w, r, CreateResponse(channels))
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
		ctx, &pb_channel.DeleteRequest{
			Id: cId,
		},
	)

	if err != nil {
		HandleGrpcError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

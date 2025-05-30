package api_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"mistapi/src/api"
	pb_channel "mistapi/src/protos/v1/channel"
	"mistapi/src/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	apiUrl = "/api/v1/channels"
)

func TestCreateChannel(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:successfully_creating_channel", func(t *testing.T) {
		// ARRANGE
		channel := api.Channel{
			ID:          "1",
			Name:        "foo-channel",
			AppserverId: "1",
		}
		expected := marshallResponse(t, api.CreateResponse(channel))
		mockCreateRequest := &pb_channel.CreateRequest{Name: channel.Name, AppserverId: channel.AppserverId}
		mockCreateResponse := &pb_channel.CreateResponse{Channel: &pb_channel.Channel{
			Id:          channel.ID,
			Name:        channel.Name,
			AppserverId: channel.AppserverId,
		}}
		mockService := new(testutil.MockChannelService)
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.ChannelCreate{Name: channel.Name, AppserverId: channel.AppserverId})
		req, err := http.NewRequest("POST", apiUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.ChannelCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(testutil.MockChannelService)
		mockCreateRequest := &pb_channel.CreateRequest{Name: "foo-channel", AppserverId: "1"}
		mockResponse := &pb_channel.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.ChannelCreate{Name: "foo-channel", AppserverId: "1"})
		req, err := http.NewRequest("POST", apiUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.ChannelCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(testutil.MockChannelService)
		mockCreateRequest := &pb_channel.CreateRequest{Name: "foo-channel", AppserverId: "1"}
		mockResponse := &pb_channel.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, "invalid")
		req, err := http.NewRequest("POST", apiUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.ChannelCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestListChannelsHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	params := url.Values{}
	appserverId := "1"
	params.Add("appserver_id", appserverId)
	urlWithParams := apiUrl + "?" + params.Encode()

	r := chi.NewRouter()
	r.Get(apiUrl, api.ListChannelsHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:successfully_returns_channels", func(t *testing.T) {
		// ARRANGE
		channels := []api.Channel{
			{ID: "1", Name: "bar", AppserverId: appserverId},
			{ID: "2", Name: "bar", AppserverId: appserverId},
		}
		expected := marshallResponse(t, api.CreateResponse(channels))
		mockRequest := &pb_channel.ListServerChannelsRequest{AppserverId: appserverId}
		mockResponse := &pb_channel.ListServerChannelsResponse{}
		mockResponse.Channels = []*pb_channel.Channel{
			{Id: channels[0].ID, Name: channels[0].Name, AppserverId: channels[0].AppserverId},
			{Id: channels[1].ID, Name: channels[1].Name, AppserverId: channels[1].AppserverId},
		}

		mockService := new(testutil.MockChannelService)
		mockService.On("ListServerChannels", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", urlWithParams, nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(testutil.MockChannelService)
		mockRequest := &pb_channel.ListServerChannelsRequest{AppserverId: appserverId}
		mockResponse := &pb_channel.ListServerChannelsResponse{}
		mockService.On("ListServerChannels", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", urlWithParams, nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:errors_when_no_appserver_id_provided", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Appserver ID is required"))
		mockService := new(testutil.MockChannelService)
		mockRequest := &pb_channel.ListServerChannelsRequest{AppserverId: appserverId}
		mockResponse := &pb_channel.ListServerChannelsResponse{}
		mockService.On("ListServerChannels", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", apiUrl, nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestDeleteChannel(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.ChannelDeleteHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:is_successful", func(t *testing.T) {
		// ARRANGE
		cId := "1"
		mockDeleteRequest := &pb_channel.DeleteRequest{Id: cId}
		mockDeleteResponse := &pb_channel.DeleteResponse{}

		mockService := new(testutil.MockChannelService)
		mockService.On(
			"Delete", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", cId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", cId)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:on_error_when_deleting_returns_error", func(t *testing.T) {
		// ARRANGE
		cId := "1"
		mockService := new(testutil.MockChannelService)
		mockDeleteRequest := &pb_channel.DeleteRequest{Id: cId}
		mockResponse := &pb_channel.DeleteResponse{}
		mockService.On("Delete", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", cId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", cId)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

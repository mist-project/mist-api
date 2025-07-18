package api_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mistapi/src/api"
	"mistapi/src/protos/v1/channel"
	"mistapi/src/testutil"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	apiUrl = "/api/v1/channels"
)

func TestCreateChannel(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:successfully_creating_channel", func(t *testing.T) {
		// ARRANGE
		c := types.Channel{
			ID:          "1",
			Name:        "foo-channel",
			AppserverId: "1",
		}
		expected := marshallResponse(t, api.CreateResponse(c))
		mockCreateRequest := &channel.CreateRequest{Name: c.Name, AppserverId: c.AppserverId}
		mockCreateResponse := &channel.CreateResponse{Channel: &channel.Channel{
			Id:          c.ID,
			Name:        c.Name,
			AppserverId: c.AppserverId,
		}}
		mockService := new(testutil.MockChannelService)
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, types.ChannelCreate{Name: c.Name, AppserverId: c.AppserverId})
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
		mockCreateRequest := &channel.CreateRequest{Name: "foo-channel", AppserverId: "1"}
		mockResponse := &channel.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, types.ChannelCreate{Name: "foo-channel", AppserverId: "1"})
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
		mockCreateRequest := &channel.CreateRequest{Name: "foo-channel", AppserverId: "1"}
		mockResponse := &channel.CreateResponse{}
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

func TestDeleteChannel(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.ChannelDeleteHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:is_successful", func(t *testing.T) {
		// ARRANGE
		cId := "1"
		mockDeleteRequest := &channel.DeleteRequest{Id: cId}
		mockDeleteResponse := &channel.DeleteResponse{}

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
		mockDeleteRequest := &channel.DeleteRequest{Id: cId}
		mockResponse := &channel.DeleteResponse{}
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

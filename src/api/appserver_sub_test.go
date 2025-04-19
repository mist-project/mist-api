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
	pb "mistapi/src/protos/v1/gen"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAppserverSub(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	url := "/api/v1/appserver-sub"

	t.Run("successfully_creating_appserver_sub", func(t *testing.T) {
		// ARRANGE
		sub := api.AppserverSub{
			ID:          "1",
			AppserverId: "1",
			AppuserId:   "",
		}
		expected := marshallResponse(t, api.CreateResponse(sub))
		mockCreateRequest := &pb.CreateAppserverSubRequest{AppserverId: sub.AppserverId}
		mockCreateResponse := &pb.CreateAppserverSubResponse{AppserverSub: &pb.AppserverSub{
			Id:          sub.ID,
			AppserverId: sub.AppserverId,
		}}
		mockService := new(MockService)
		mockService.On("CreateAppserverSub", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverSubCreate{AppserverId: sub.AppserverId})
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverSubCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverSubRequest{AppserverId: "1"}
		mockResponse := &pb.CreateAppserverSubResponse{}
		mockService.On("CreateAppserverSub", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverSubCreate{AppserverId: "1"})
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverSubCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverSubRequest{AppserverId: "1"}
		mockResponse := &pb.CreateAppserverSubResponse{}
		mockService.On("CreateAppserverSub", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, "invalid")
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverSubCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestDeleteAppserverSub(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.AppserverSubDeleteHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("is_successful", func(t *testing.T) {
		// ARRANGE
		sId := "1"
		mockDeleteRequest := &pb.DeleteAppserverSubRequest{Id: sId}
		mockDeleteResponse := &pb.DeleteAppserverSubResponse{}

		mockService := new(MockService)
		mockService.On(
			"DeleteAppserverSub", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", sId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", sId)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("on_error_when_deleting_returns_error", func(t *testing.T) {
		// ARRANGE
		sId := "1"
		mockService := new(MockService)
		mockDeleteRequest := &pb.DeleteAppserverSubRequest{Id: sId}
		mockResponse := &pb.DeleteAppserverSubResponse{}
		mockService.On("DeleteAppserverSub", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", sId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", sId)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

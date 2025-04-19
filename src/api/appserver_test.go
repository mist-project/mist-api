package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mistapi/src/api"
	pb "mistapi/src/protos/v1/gen"
)

func TestList(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("successfully_returns_appservers", func(t *testing.T) {
		// ARRANGE
		servers := []api.Appserver{
			{ID: "1", Name: "bar", IsOwner: true},
			{ID: "2", Name: "bar", IsOwner: true},
		}

		mockResponse := &pb.ListAppserversResponse{}
		mockResponse.Appservers = []*pb.Appserver{
			{Id: servers[0].ID, Name: servers[0].Name, IsOwner: servers[0].IsOwner},
			{Id: servers[1].ID, Name: servers[1].Name, IsOwner: servers[1].IsOwner},
		}

		mockService := new(MockService)
		mockService.On("ListAppservers", mock.Anything, mock.Anything).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appserver", nil)
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.List(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)

		expected, _ := json.Marshal(api.CreateResponse(servers))
		assert.JSONEq(t, string(expected), rr.Body.String())
	})

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		mockService := new(MockService)
		mockResponse := &pb.ListAppserversResponse{}
		mockService.On("ListAppservers", mock.Anything, mock.Anything).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appserver", nil)
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.List(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expected, _ := json.Marshal(api.CreateErrorResponse("Bad request"))
		assert.JSONEq(t, string(expected), rr.Body.String())
	})
}

func TestListSubs(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("successfully_returns_appservers_and_sub_id", func(t *testing.T) {

		// ARRANGE
		servers := []api.AppserverAndSub{
			{Appserver: api.Appserver{ID: "1", Name: "bar", IsOwner: true}, SubId: "1"},
			{Appserver: api.Appserver{ID: "2", Name: "bar", IsOwner: true}, SubId: "2"},
		}

		mockResponse := &pb.GetUserAppserverSubsResponse{}
		mockResponse.Appservers = []*pb.AppserverAndSub{
			{Appserver: &pb.Appserver{
				Id:      servers[0].Appserver.ID,
				Name:    servers[0].Appserver.Name,
				IsOwner: servers[0].Appserver.IsOwner},
				SubId: servers[0].SubId},
			{Appserver: &pb.Appserver{
				Id:      servers[1].Appserver.ID,
				Name:    servers[1].Appserver.Name,
				IsOwner: servers[1].Appserver.IsOwner},
				SubId: servers[1].SubId},
		}

		mockService := new(MockService)
		mockService.On("GetUserAppserverSubs", mock.Anything, mock.Anything).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appserver/subs", nil)
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.ListSubs(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)

		jsonData, _ := json.Marshal(api.CreateResponse(servers))
		assert.JSONEq(t, string(jsonData), rr.Body.String())
	})

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		mockService := new(MockService)
		mockResponse := &pb.GetUserAppserverSubsResponse{}
		mockService.On("GetUserAppserverSubs", mock.Anything, mock.Anything).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appserver/subs", nil)
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.ListSubs(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		expected, _ := json.Marshal(api.CreateErrorResponse("Bad request"))
		assert.JSONEq(t, string(expected), rr.Body.String())
	})
}

func TestCreateAppserver(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("successfully_creating_appserver", func(t *testing.T) {
		// ARRANGE
		appserver := &pb.Appserver{
			Id:      "1",
			Name:    "foo",
			IsOwner: true,
		}
		mockCreateRequest := &pb.CreateAppserverRequest{Name: appserver.Name}
		mockCreateResponse := &pb.CreateAppserverResponse{Appserver: appserver}
		mockService := new(MockService)
		mockService.On("CreateAppserver", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		jsonData, err := json.Marshal(api.AppserverCreate{Name: appserver.Name})

		req, err := http.NewRequest("POST", "/api/v1/appserver", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.Create(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusCreated, rr.Code)

		expected, _ := json.Marshal(api.CreateResponse(appserver))
		assert.JSONEq(t, string(expected), rr.Body.String())
	})

	t.Run("errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverRequest{Name: "foo"}
		mockResponse := &pb.CreateAppserverResponse{}
		mockService.On("CreateAppserver", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		jsonData, err := json.Marshal(api.AppserverCreate{Name: "foo"})

		req, err := http.NewRequest("POST", "/api/v1/appserver", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.Create(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expected, _ := json.Marshal(api.CreateErrorResponse("Internal Server Error."))
		assert.JSONEq(t, string(expected), rr.Body.String())
	})

	t.Run("errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverRequest{Name: "foo"}
		mockResponse := &pb.CreateAppserverResponse{}
		mockService.On("CreateAppserver", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		jsonData, err := json.Marshal("ok")

		req, err := http.NewRequest("POST", "/api/v1/appserver", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.Create(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		expected, _ := json.Marshal(api.CreateErrorResponse("Invalid attributes provided."))
		assert.JSONEq(t, string(expected), rr.Body.String())
	})
}

func TestDeleteAppserver(t *testing.T) {

	// add router
	r := chi.NewRouter()
	r.Delete("/{id}", api.Delete)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("is_successful", func(t *testing.T) {
		// ARRANGE
		someid := "someid"
		mockDeleteRequest := &pb.DeleteAppserverRequest{Id: someid}
		mockDeleteResponse := &pb.DeleteAppserverResponse{}

		mockService := new(MockService)
		mockService.On(
			"DeleteAppserver", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", "/someid", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", someid)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("on_error_when_deleting_returns_error", func(t *testing.T) {
		// ARRANGE
		someid := "someid"
		mockService := new(MockService)
		mockDeleteRequest := &pb.DeleteAppserverRequest{Id: someid}
		mockResponse := &pb.DeleteAppserverResponse{}
		mockService.On("DeleteAppserver", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", "/someid", nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", someid)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

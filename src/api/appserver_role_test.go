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

func TestCreateAppserverRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	url := "/api/v1/appserver-role"

	t.Run("successfully_creating_appserver_role", func(t *testing.T) {
		// ARRANGE
		role := &pb.AppserverRole{
			Id:          "1",
			Name:        "foo",
			AppserverId: "1",
		}
		expected := marshallResponse(t, api.CreateResponse(role))
		mockCreateRequest := &pb.CreateAppserverRoleRequest{Name: role.Name, AppserverId: role.AppserverId}
		mockCreateResponse := &pb.CreateAppserverRoleResponse{AppserverRole: role}
		mockService := new(MockService)
		mockService.On("CreateAppserverRole", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverRoleCreate{Name: role.Name, AppserverId: role.AppserverId})
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverRoleCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverRoleRequest{Name: "foo", AppserverId: "1"}
		mockResponse := &pb.CreateAppserverRoleResponse{}
		mockService.On("CreateAppserverRole", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverRoleCreate{Name: "foo", AppserverId: "1"})
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverRoleCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverRoleRequest{Name: "foo", AppserverId: "1"}
		mockResponse := &pb.CreateAppserverRoleResponse{}
		mockService.On("CreateAppserverRole", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

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
		api.AppserverRoleCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestDeleteAppserverRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.AppserverRoleDeleteHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("is_successful", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockDeleteRequest := &pb.DeleteAppserverRoleRequest{Id: aId}
		mockDeleteResponse := &pb.DeleteAppserverRoleResponse{}

		mockService := new(MockService)
		mockService.On(
			"DeleteAppserverRole", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", aId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", aId)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("on_error_when_deleting_returns_error", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockService := new(MockService)
		mockDeleteRequest := &pb.DeleteAppserverRoleRequest{Id: aId}
		mockResponse := &pb.DeleteAppserverRoleResponse{}
		mockService.On("DeleteAppserverRole", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", aId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", aId)

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

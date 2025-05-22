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
	pb_appserver_role "mistapi/src/protos/v1/appserver_role"
	"mistapi/src/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAppserverRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	url := "/api/v1/appserver-roles"

	t.Run("Success:successfully_creating_appserver_role", func(t *testing.T) {
		// ARRANGE
		role := api.AppserverRole{
			ID:          "1",
			Name:        "foo",
			AppserverId: "1",
		}
		expected := marshallResponse(t, api.CreateResponse(role))
		mockCreateRequest := &pb_appserver_role.CreateRequest{Name: role.Name, AppserverId: role.AppserverId}
		mockCreateResponse := &pb_appserver_role.CreateResponse{AppserverRole: &pb_appserver_role.AppserverRole{
			Id:          role.ID,
			Name:        role.Name,
			AppserverId: role.AppserverId,
		}}
		mockService := new(testutil.MockAppserverRoleService)
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

	t.Run("Error:errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(testutil.MockAppserverRoleService)
		mockCreateRequest := &pb_appserver_role.CreateRequest{Name: "foo", AppserverId: "1"}
		mockResponse := &pb_appserver_role.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

	t.Run("Error:errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(testutil.MockAppserverRoleService)
		mockCreateRequest := &pb_appserver_role.CreateRequest{Name: "foo", AppserverId: "1"}
		mockResponse := &pb_appserver_role.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

	t.Run("Success:is_successful", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockDeleteRequest := &pb_appserver_role.DeleteRequest{Id: aId}
		mockDeleteResponse := &pb_appserver_role.DeleteResponse{}

		mockService := new(testutil.MockAppserverRoleService)
		mockService.On(
			"Delete", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

	t.Run("Error:on_error_when_deleting_returns_error", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockService := new(testutil.MockAppserverRoleService)
		mockDeleteRequest := &pb_appserver_role.DeleteRequest{Id: aId}
		mockResponse := &pb_appserver_role.DeleteResponse{}
		mockService.On("Delete", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

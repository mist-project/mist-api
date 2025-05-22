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
	pb_appserver_permission "mistapi/src/protos/v1/appserver_permission"
	"mistapi/src/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAppserverPermission(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	url := "/api/v1/appserver-roles"

	t.Run("Success:successfully_creating_appserver_role", func(t *testing.T) {
		// ARRANGE
		role := api.AppserverPermission{
			ID:          "1",
			AppuserId:   "foo",
			AppserverId: "1",
		}
		mockCreateRequest := &pb_appserver_permission.CreateRequest{AppuserId: role.AppuserId, AppserverId: role.AppserverId}
		mockCreateResponse := &pb_appserver_permission.CreateResponse{}
		mockService := new(testutil.MockAppserverPermissionService)
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverPermissionCreate{AppuserId: role.AppuserId, AppserverId: role.AppserverId})
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverPermissionCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(testutil.MockAppserverPermissionService)
		mockCreateRequest := &pb_appserver_permission.CreateRequest{AppuserId: "boom", AppserverId: "boom"}
		mockResponse := &pb_appserver_permission.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverPermissionCreate{AppuserId: "boom", AppserverId: "boom"})
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverPermissionCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(testutil.MockAppserverPermissionService)
		mockCreateRequest := &pb_appserver_permission.CreateRequest{AppuserId: "boom", AppserverId: "boom"}
		mockResponse := &pb_appserver_permission.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, "invalid")
		req, err := http.NewRequest("POST", url, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverPermissionCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestListAppserverPermissionUsersHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	params := url.Values{}
	appserverId := "1"
	params.Add("appserver_id", appserverId)
	urlWithParams := apiUrl + "?" + params.Encode()

	r := chi.NewRouter()
	r.Get(apiUrl, api.AppserverPermissionListHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:successfully_returns_appserver_permissions", func(t *testing.T) {
		// ARRANGE
		permissions := []api.AppserverPermission{
			{ID: "1", AppuserId: "user", AppserverId: appserverId},
			{ID: "2", AppuserId: "user", AppserverId: appserverId},
		}
		expected := marshallResponse(t, api.CreateResponse(permissions))
		mockRequest := &pb_appserver_permission.ListAppserverUsersRequest{AppserverId: appserverId}
		mockResponse := &pb_appserver_permission.ListAppserverUsersResponse{}
		mockResponse.AppserverPermissions = []*pb_appserver_permission.AppserverPermission{
			{Id: permissions[0].ID, AppuserId: permissions[0].AppuserId, AppserverId: permissions[0].AppserverId},
			{Id: permissions[1].ID, AppuserId: permissions[1].AppuserId, AppserverId: permissions[1].AppserverId},
		}

		mockService := new(testutil.MockAppserverPermissionService)
		mockService.On("ListAppserverUsers", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
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
		mockService := new(testutil.MockAppserverPermissionService)
		mockRequest := &pb_appserver_permission.ListAppserverUsersRequest{AppserverId: appserverId}
		mockResponse := &pb_appserver_permission.ListAppserverUsersResponse{}
		mockService.On("ListAppserverUsers", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
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
		mockService := new(testutil.MockAppserverPermissionService)
		mockRequest := &pb_appserver_permission.ListAppserverUsersRequest{AppserverId: appserverId}
		mockResponse := &pb_appserver_permission.ListAppserverUsersResponse{}
		mockService.On("ListAppserverUsers", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
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

func TestDeleteAppserverPermission(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.AppserverPermissionDeleteHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:is_successful", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockDeleteRequest := &pb_appserver_permission.DeleteRequest{Id: aId}
		mockDeleteResponse := &pb_appserver_permission.DeleteResponse{}

		mockService := new(testutil.MockAppserverPermissionService)
		mockService.On(
			"Delete", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
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
		mockService := new(testutil.MockAppserverPermissionService)
		mockDeleteRequest := &pb_appserver_permission.DeleteRequest{Id: aId}
		mockResponse := &pb_appserver_permission.DeleteResponse{}
		mockService.On("Delete", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverPermissionClient").Return(mockService)
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

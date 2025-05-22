package api_test

import (
	"errors"
	"fmt"
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
	pb_appserver "mistapi/src/protos/v1/appserver"
	pb_appserver_role "mistapi/src/protos/v1/appserver_role"
	pb_appserver_sub "mistapi/src/protos/v1/appserver_sub"
	pb_appuser "mistapi/src/protos/v1/appuser"
	"mistapi/src/testutil"
)

func TestListAppservers(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:successfully_returns_appservers_and_sub_id", func(t *testing.T) {
		// ARRANGE
		servers := []api.AppserverAndSub{
			{Appserver: api.Appserver{ID: "1", Name: "bar", IsOwner: true}, SubId: "1"},
			{Appserver: api.Appserver{ID: "1", Name: "bar", IsOwner: true}, SubId: "1"},
		}
		expected := marshallResponse(t, api.CreateResponse(servers))
		mockRequest := &pb_appserver_sub.ListUserServerSubsRequest{}
		mockResponse := &pb_appserver_sub.ListUserServerSubsResponse{}
		mockResponse.Appservers = []*pb_appserver_sub.AppserverAndSub{
			{
				Appserver: &pb_appserver.Appserver{
					Id:      servers[0].Appserver.ID,
					Name:    servers[0].Appserver.Name,
					IsOwner: servers[0].Appserver.IsOwner},
				SubId: servers[0].SubId,
			},
			{
				Appserver: &pb_appserver.Appserver{
					Id:      servers[1].Appserver.ID,
					Name:    servers[1].Appserver.Name,
					IsOwner: servers[1].Appserver.IsOwner},
				SubId: servers[1].SubId,
			},
		}

		mockService := new(testutil.MockAppserverSubService)
		mockService.On("ListUserServerSubs", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverSubClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appservers", nil)
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverListHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(testutil.MockAppserverSubService)
		mockRequest := &pb_appserver_sub.ListUserServerSubsRequest{}
		mockResponse := &pb_appserver_sub.ListUserServerSubsResponse{}
		mockService.On("ListUserServerSubs", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverSubClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appservers", nil)
		require.NoError(t, err)

		// Mock the token in the request context
		req = addContextHeaders(req)

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverListHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestCreateAppserver(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:successfully_creating_appserver", func(t *testing.T) {
		// ARRANGE
		appserver := &pb_appserver.Appserver{
			Id:      "1",
			Name:    "foo",
			IsOwner: true,
		}
		expected := marshallResponse(t, api.CreateResponse(appserver))
		mockCreateRequest := &pb_appserver.CreateRequest{Name: appserver.Name}
		mockCreateResponse := &pb_appserver.CreateResponse{Appserver: appserver}
		mockService := new(testutil.MockAppserverService)
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverCreate{Name: appserver.Name})
		req, err := http.NewRequest("POST", "/api/v1/appservers", payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(testutil.MockAppserverService)
		mockCreateRequest := &pb_appserver.CreateRequest{Name: "foo"}
		mockResponse := &pb_appserver.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, api.AppserverCreate{Name: "foo"})
		req, err := http.NewRequest("POST", "/api/v1/appservers", payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(testutil.MockAppserverService)
		mockCreateRequest := &pb_appserver.CreateRequest{Name: "foo"}
		mockResponse := &pb_appserver.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, "invalid")
		req, err := http.NewRequest("POST", "/api/v1/appservers", payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverCreateHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestDeleteAppserver(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.AppserverDeleteHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:is_successful", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockDeleteRequest := &pb_appserver.DeleteRequest{Id: aId}
		mockDeleteResponse := &pb_appserver.DeleteResponse{}

		mockService := new(testutil.MockAppserverService)
		mockService.On(
			"Delete", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
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
		mockService := new(testutil.MockAppserverService)
		mockDeleteRequest := &pb_appserver.DeleteRequest{Id: aId}
		mockResponse := &pb_appserver.DeleteResponse{}
		mockService.On("Delete", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
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

func TestAppserverDetailHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverDetailHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:successfully_returns_appserver_details", func(t *testing.T) {

		// ARRANGE
		appserver := api.AppserverDetail{
			ID:      "1",
			Name:    "Foo",
			IsOwner: false,
		}
		expected := marshallResponse(t, api.CreateResponse(appserver))

		mockRequest := &pb_appserver.GetByIdRequest{Id: appserver.ID}
		mockResponse := &pb_appserver.GetByIdResponse{Appserver: &pb_appserver.Appserver{
			Id:      appserver.ID,
			Name:    appserver.Name,
			IsOwner: appserver.IsOwner,
		}}

		mockService := new(testutil.MockAppserverService)
		mockService.On("GetById", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", appserver.ID), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", appserver.ID)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		appserver := api.AppserverDetail{
			ID: "1",
		}
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockRequest := &pb_appserver.GetByIdRequest{Id: appserver.ID}
		mockResponse := &pb_appserver.GetByIdResponse{}

		mockService := new(testutil.MockAppserverService)
		mockService.On("GetById", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", appserver.ID), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", appserver.ID)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestAppserverListSubsHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverListSubsHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	aId := "123"

	t.Run("Success:successfully_returns_appuser_and_sub_ids", func(t *testing.T) {

		// ARRANGE
		appusers := []api.AppuserAppserverSub{
			{Appuser: api.Appuser{ID: "1", Username: "foo"}, SubId: "1"},
			{Appuser: api.Appuser{ID: "2", Username: "bar"}, SubId: "2"},
		}
		expected := marshallResponse(t, api.CreateResponse(appusers))
		mockResponse := &pb_appserver_sub.ListAppserverUserSubsResponse{}
		mockResponse.Appusers = []*pb_appserver_sub.AppuserAndSub{
			{Appuser: &pb_appuser.Appuser{
				Id:       appusers[0].Appuser.ID,
				Username: appusers[0].Appuser.Username,
			},
				SubId: appusers[0].SubId},
			{Appuser: &pb_appuser.Appuser{
				Id:       appusers[1].Appuser.ID,
				Username: appusers[1].Appuser.Username,
			},
				SubId: appusers[1].SubId},
		}

		mockRequest := &pb_appserver_sub.ListAppserverUserSubsRequest{AppserverId: aId}

		mockSubService := new(testutil.MockAppserverSubService)
		mockSubService.On("ListAppserverUserSubs", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverSubClient").Return(mockSubService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", aId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", aId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockRequest := &pb_appserver_sub.ListAppserverUserSubsRequest{AppserverId: aId}
		mockSubService := new(testutil.MockAppserverSubService)
		mockSubService.On(
			"ListAppserverUserSubs", mock.Anything, mockRequest,
		).Return(nil, status.Error(codes.InvalidArgument, "Bad request"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverSubClient").Return(mockSubService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", aId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", aId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestAppserverListRolesHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverListRolesHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	aId := "123"

	t.Run("Success:successfully_returns_appserver_roles", func(t *testing.T) {

		// ARRANGE
		roles := []api.AppserverRole{
			{ID: "1", Name: "foo", AppserverId: aId},
			{ID: "2", Name: "bar", AppserverId: aId},
		}
		expected := marshallResponse(t, api.CreateResponse(roles))
		mockRequest := &pb_appserver_role.ListServerRolesRequest{AppserverId: aId}
		mockResponse := &pb_appserver_role.ListServerRolesResponse{}
		mockResponse.AppserverRoles = []*pb_appserver_role.AppserverRole{
			{Id: roles[0].ID, Name: roles[0].Name, AppserverId: roles[0].AppserverId},
			{Id: roles[1].ID, Name: roles[1].Name, AppserverId: roles[1].AppserverId},
		}

		// mockService := new(testutil.MockAppserverService)
		// mockService.On("ListServerRoles", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockRoleService := new(testutil.MockAppserverRoleService)
		mockRoleService.On("ListServerRoles", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		// mockClient.On("GetAppserverClient").Return(mockService)
		mockClient.On("GetAppserverRoleClient").Return(mockRoleService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", aId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", aId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(testutil.MockAppserverRoleService)
		mockResponse := &pb_appserver_role.ListServerRolesResponse{}
		mockRequest := &pb_appserver_role.ListServerRolesRequest{AppserverId: aId}
		mockService.On("ListServerRoles", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", aId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", aId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

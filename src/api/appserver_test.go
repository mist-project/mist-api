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
		expected := marshallResponse(t, api.CreateResponse(servers))
		mockRequest := &pb.ListAppserversRequest{}
		mockResponse := &pb.ListAppserversResponse{}
		mockResponse.Appservers = []*pb.Appserver{
			{Id: servers[0].ID, Name: servers[0].Name, IsOwner: servers[0].IsOwner},
			{Id: servers[1].ID, Name: servers[1].Name, IsOwner: servers[1].IsOwner},
		}

		mockService := new(MockService)
		mockService.On("ListAppservers", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(MockService)
		mockRequest := &pb.ListAppserversRequest{}
		mockResponse := &pb.ListAppserversResponse{}
		mockService.On("ListAppservers", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

func TestAppserverUserListSubsHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("successfully_returns_appservers_and_sub_id", func(t *testing.T) {

		// ARRANGE
		servers := []api.AppserverAndSub{
			{Appserver: api.Appserver{ID: "1", Name: "bar", IsOwner: true}, SubId: "1"},
			{Appserver: api.Appserver{ID: "2", Name: "bar", IsOwner: true}, SubId: "2"},
		}
		expected := marshallResponse(t, api.CreateResponse(servers))
		mockRequest := &pb.GetUserAppserverSubsRequest{}
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
		mockService.On("GetUserAppserverSubs", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appservers/subs", nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverUserListSubsHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(MockService)
		mockRequest := &pb.GetUserAppserverSubsRequest{}
		mockResponse := &pb.GetUserAppserverSubsResponse{}
		mockService.On("GetUserAppserverSubs", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", "/api/v1/appservers/subs", nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverUserListSubsHandler(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
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
		expected := marshallResponse(t, api.CreateResponse(appserver))
		mockCreateRequest := &pb.CreateAppserverRequest{Name: appserver.Name}
		mockCreateResponse := &pb.CreateAppserverResponse{Appserver: appserver}
		mockService := new(MockService)
		mockService.On("CreateAppserver", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("errors_during_creation_returns_error_status", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Internal Server Error."))
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverRequest{Name: "foo"}
		mockResponse := &pb.CreateAppserverResponse{}
		mockService.On("CreateAppserver", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("errors_with_invalid_post_parameters", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))
		mockService := new(MockService)
		mockCreateRequest := &pb.CreateAppserverRequest{Name: "foo"}
		mockResponse := &pb.CreateAppserverResponse{}
		mockService.On("CreateAppserver", mock.Anything, mockCreateRequest).Return(mockResponse, nil)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("is_successful", func(t *testing.T) {
		// ARRANGE
		aId := "1"
		mockDeleteRequest := &pb.DeleteAppserverRequest{Id: aId}
		mockDeleteResponse := &pb.DeleteAppserverResponse{}

		mockService := new(MockService)
		mockService.On(
			"DeleteAppserver", mock.Anything, mockDeleteRequest,
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
		mockDeleteRequest := &pb.DeleteAppserverRequest{Id: aId}
		mockResponse := &pb.DeleteAppserverResponse{}
		mockService.On("DeleteAppserver", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
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

func TestAppserverDetailHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverDetailHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("successfully_returns_appserver_details", func(t *testing.T) {

		// ARRANGE
		appserver := api.AppserverDetail{
			ID:      "1",
			Name:    "Foo",
			IsOwner: false,
		}
		expected := marshallResponse(t, api.CreateResponse(appserver))

		mockRequest := &pb.GetByIdAppserverRequest{Id: appserver.ID}
		mockResponse := &pb.GetByIdAppserverResponse{Appserver: &pb.Appserver{
			Id:      appserver.ID,
			Name:    appserver.Name,
			IsOwner: appserver.IsOwner,
		}}

		mockService := new(MockService)
		mockService.On("GetByIdAppserver", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		appserver := api.AppserverDetail{
			ID: "1",
		}
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockRequest := &pb.GetByIdAppserverRequest{Id: appserver.ID}
		mockResponse := &pb.GetByIdAppserverResponse{}

		mockService := new(MockService)
		mockService.On("GetByIdAppserver", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("successfully_returns_appuser_and_sub_ids", func(t *testing.T) {

		// ARRANGE
		appusers := []api.AppuserAppserverSub{
			{Appuser: api.Appuser{ID: "1", Username: "foo"}, SubId: "1"},
			{Appuser: api.Appuser{ID: "2", Username: "bar"}, SubId: "2"},
		}
		expected := marshallResponse(t, api.CreateResponse(appusers))
		mockResponse := &pb.GetAllUsersAppserverSubsResponse{}
		mockResponse.Appusers = []*pb.AppuserAndSub{
			{Appuser: &pb.Appuser{
				Id:       appusers[0].Appuser.ID,
				Username: appusers[0].Appuser.Username,
			},
				SubId: appusers[0].SubId},
			{Appuser: &pb.Appuser{
				Id:       appusers[1].Appuser.ID,
				Username: appusers[1].Appuser.Username,
			},
				SubId: appusers[1].SubId},
		}

		mockRequest := &pb.GetAllUsersAppserverSubsRequest{AppserverId: aId}

		mockService := new(MockService)
		mockService.On("GetAllUsersAppserverSubs", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(MockService)
		mockResponse := &pb.GetAllUsersAppserverSubsResponse{}
		mockRequest := &pb.GetAllUsersAppserverSubsRequest{AppserverId: aId}
		mockService.On("GetAllUsersAppserverSubs", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("successfully_returns_appserver_roles", func(t *testing.T) {

		// ARRANGE
		roles := []api.AppserverRole{
			{ID: "1", Name: "foo", AppserverId: aId},
			{ID: "2", Name: "bar", AppserverId: aId},
		}
		expected := marshallResponse(t, api.CreateResponse(roles))
		mockRequest := &pb.GetAllAppserverRolesRequest{AppserverId: aId}
		mockResponse := &pb.GetAllAppserverRolesResponse{}
		mockResponse.AppserverRoles = []*pb.AppserverRole{
			{Id: roles[0].ID, Name: roles[0].Name, AppserverId: roles[0].AppserverId},
			{Id: roles[1].ID, Name: roles[1].Name, AppserverId: roles[1].AppserverId},
		}

		mockService := new(MockService)
		mockService.On("GetAllAppserverRoles", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

	t.Run("on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockService := new(MockService)
		mockResponse := &pb.GetAllAppserverRolesResponse{}
		mockRequest := &pb.GetAllAppserverRolesRequest{AppserverId: aId}
		mockService.On("GetAllAppserverRoles", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)

		mockClient := new(MockClient)
		mockClient.On("GetServerClient").Return(mockService)
		MockGrpcClient(t, mockClient)

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

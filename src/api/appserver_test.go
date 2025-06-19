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
	"mistapi/src/protos/v1/appserver"
	"mistapi/src/protos/v1/appserver_role"
	"mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/protos/v1/appserver_sub"
	"mistapi/src/protos/v1/appuser"
	"mistapi/src/protos/v1/channel"
	"mistapi/src/protos/v1/channel_role"
	"mistapi/src/testutil"
	"mistapi/src/types"
)

func TestListAppservers(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:successfully_returns_appservers_and_sub_id", func(t *testing.T) {
		// ARRANGE
		servers := []types.AppserverAndSub{
			{Appserver: types.Appserver{ID: "1", Name: "bar", IsOwner: true}, SubId: "1"},
			{Appserver: types.Appserver{ID: "1", Name: "bar", IsOwner: true}, SubId: "1"},
		}
		expected := marshallResponse(t, api.CreateResponse(servers))
		mockRequest := &appserver_sub.ListUserServerSubsRequest{}
		mockResponse := &appserver_sub.ListUserServerSubsResponse{}
		mockResponse.Appservers = []*appserver_sub.AppserverAndSub{
			{
				Appserver: &appserver.Appserver{
					Id:      servers[0].Appserver.ID,
					Name:    servers[0].Appserver.Name,
					IsOwner: servers[0].Appserver.IsOwner},
				SubId: servers[0].SubId,
			},
			{
				Appserver: &appserver.Appserver{
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
		mockRequest := &appserver_sub.ListUserServerSubsRequest{}
		mockResponse := &appserver_sub.ListUserServerSubsResponse{}
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
		s := &appserver.Appserver{
			Id:      "1",
			Name:    "foo",
			IsOwner: true,
		}
		expected := marshallResponse(t, api.CreateResponse(s))
		mockCreateRequest := &appserver.CreateRequest{Name: s.Name}
		mockCreateResponse := &appserver.CreateResponse{Appserver: s}
		mockService := new(testutil.MockAppserverService)
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockCreateResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, types.AppserverCreate{Name: s.Name})
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
		mockCreateRequest := &appserver.CreateRequest{Name: "foo"}
		mockResponse := &appserver.CreateResponse{}
		mockService.On("Create", mock.Anything, mockCreateRequest).Return(mockResponse, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		payload := marshallPayload(t, types.AppserverCreate{Name: "foo"})
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
		mockCreateRequest := &appserver.CreateRequest{Name: "foo"}
		mockResponse := &appserver.CreateResponse{}
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
		sId := "1"
		mockDeleteRequest := &appserver.DeleteRequest{Id: sId}
		mockDeleteResponse := &appserver.DeleteResponse{}

		mockService := new(testutil.MockAppserverService)
		mockService.On(
			"Delete", mock.Anything, mockDeleteRequest,
		).Return(mockDeleteResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

	t.Run("Error:on_error_when_deleting_returns_error", func(t *testing.T) {
		// ARRANGE
		sId := "1"
		mockService := new(testutil.MockAppserverService)
		mockDeleteRequest := &appserver.DeleteRequest{Id: sId}
		mockResponse := &appserver.DeleteResponse{}
		mockService.On("Delete", mock.Anything, mockDeleteRequest).Return(mockResponse, errors.New("boom"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

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

func TestAppserverDetailHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverDetailHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Success:successfully_returns_appserver_details", func(t *testing.T) {

		// ARRANGE
		s := types.AppserverDetail{
			ID:      "1",
			Name:    "Foo",
			IsOwner: false,
		}
		expected := marshallResponse(t, api.CreateResponse(s))

		mockRequest := &appserver.GetByIdRequest{Id: s.ID}
		mockResponse := &appserver.GetByIdResponse{Appserver: &appserver.Appserver{
			Id:      s.ID,
			Name:    s.Name,
			IsOwner: s.IsOwner,
		}}

		mockService := new(testutil.MockAppserverService)
		mockService.On("GetById", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", s.ID), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", s.ID)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		s := types.AppserverDetail{
			ID: "1",
		}
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockRequest := &appserver.GetByIdRequest{Id: s.ID}
		mockResponse := &appserver.GetByIdResponse{}

		mockService := new(testutil.MockAppserverService)
		mockService.On("GetById", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", s.ID), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", s.ID)

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

	sId := "123"

	t.Run("Success:successfully_returns_appuser_and_sub_ids", func(t *testing.T) {

		// ARRANGE
		appusers := []types.AppuserAppserverSub{
			{Appuser: types.Appuser{ID: "1", Username: "foo"}, SubId: "1"},
			{Appuser: types.Appuser{ID: "2", Username: "bar"}, SubId: "2"},
		}
		expected := marshallResponse(t, api.CreateResponse(appusers))
		mockResponse := &appserver_sub.ListAppserverUserSubsResponse{}
		mockResponse.Appusers = []*appserver_sub.AppuserAndSub{
			{Appuser: &appuser.Appuser{
				Id:       appusers[0].Appuser.ID,
				Username: appusers[0].Appuser.Username,
			},
				SubId: appusers[0].SubId},
			{Appuser: &appuser.Appuser{
				Id:       appusers[1].Appuser.ID,
				Username: appusers[1].Appuser.Username,
			},
				SubId: appusers[1].SubId},
		}

		mockRequest := &appserver_sub.ListAppserverUserSubsRequest{AppserverId: sId}

		mockSubService := new(testutil.MockAppserverSubService)
		mockSubService.On("ListAppserverUserSubs", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverSubClient").Return(mockSubService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", sId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", sId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:on_error_returns_error", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Bad request"))
		mockRequest := &appserver_sub.ListAppserverUserSubsRequest{AppserverId: sId}
		mockSubService := new(testutil.MockAppserverSubService)
		mockSubService.On(
			"ListAppserverUserSubs", mock.Anything, mockRequest,
		).Return(nil, status.Error(codes.InvalidArgument, "Bad request"))
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverSubClient").Return(mockSubService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", sId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", sId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestAppserverRoleSubListHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	sId := "123"

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverListRoleSubHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	fullUrl := fmt.Sprintf("/%s", sId)

	t.Run("Success:lists_role_subs", func(t *testing.T) {
		mockService := new(testutil.MockAppserverRoleSubService)
		mockRequest := &appserver_role_sub.ListServerRoleSubsRequest{AppserverId: sId}
		mockResponse := &appserver_role_sub.ListServerRoleSubsResponse{
			AppserverRoleSubs: []*appserver_role_sub.AppserverRoleSub{
				{Id: "1", AppuserId: "user", AppserverRoleId: "role", AppserverId: sId},
			},
		}
		mockService.On("ListServerRoleSubs", mock.Anything, mockRequest).Return(mockResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fullUrl, nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		expected := marshallResponse(t, api.CreateResponse([]types.AppserverRoleSub{
			{ID: "1", AppuserId: "user", AppserverRoleId: "role", AppserverId: sId},
		}))

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_error", func(t *testing.T) {
		mockService := new(testutil.MockAppserverRoleSubService)
		mockRequest := &appserver_role_sub.ListServerRoleSubsRequest{AppserverId: sId}
		mockService.On("ListServerRoleSubs", mock.Anything, mockRequest).Return(
			nil, status.Error(codes.InvalidArgument, "bad"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fullUrl, nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		expected := marshallResponse(t, api.CreateErrorResponse("bad"))
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestListChannelsHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	sId := "123"

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverListChannelsHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	fullUrl := fmt.Sprintf("/%s", sId)

	t.Run("Success:successfully_returns_channels", func(t *testing.T) {
		// ARRANGE
		channels := []types.Channel{
			{ID: "1", Name: "bar", AppserverId: sId},
			{ID: "2", Name: "bar", AppserverId: sId},
		}
		expected := marshallResponse(t, api.CreateResponse(channels))
		mockRequest := &channel.ListServerChannelsRequest{AppserverId: sId}
		mockResponse := &channel.ListServerChannelsResponse{}
		mockResponse.Channels = []*channel.Channel{
			{Id: channels[0].ID, Name: channels[0].Name, AppserverId: channels[0].AppserverId},
			{Id: channels[1].ID, Name: channels[1].Name, AppserverId: channels[1].AppserverId},
		}

		mockService := new(testutil.MockChannelService)
		mockService.On("ListServerChannels", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fullUrl, nil)
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
		mockRequest := &channel.ListServerChannelsRequest{AppserverId: sId}
		mockResponse := &channel.ListServerChannelsResponse{}
		mockService.On("ListServerChannels", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fullUrl, nil)
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

func TestAppserverListRolesHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Get("/{id}", api.AppserverListRolesHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	sId := "123"

	t.Run("Success:successfully_returns_appserver_roles", func(t *testing.T) {

		// ARRANGE
		roles := []types.AppserverRole{
			{ID: "1", Name: "foo", AppserverId: sId},
			{ID: "2", Name: "bar", AppserverId: sId},
		}
		expected := marshallResponse(t, api.CreateResponse(roles))
		mockRequest := &appserver_role.ListServerRolesRequest{AppserverId: sId}
		mockResponse := &appserver_role.ListServerRolesResponse{}
		mockResponse.AppserverRoles = []*appserver_role.AppserverRole{
			{Id: roles[0].ID, Name: roles[0].Name, AppserverId: roles[0].AppserverId},
			{Id: roles[1].ID, Name: roles[1].Name, AppserverId: roles[1].AppserverId},
		}

		mockRoleService := new(testutil.MockAppserverRoleService)
		mockRoleService.On("ListServerRoles", mock.Anything, mockRequest).Return(mockResponse, nil)
		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockRoleService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", sId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", sId)

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
		mockResponse := &appserver_role.ListServerRolesResponse{}
		mockRequest := &appserver_role.ListServerRolesRequest{AppserverId: sId}
		mockService.On("ListServerRoles", mock.Anything, mockRequest).Return(
			mockResponse, status.Error(codes.InvalidArgument, "Bad request"),
		)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		// Prepare the HTTP request
		req, err := http.NewRequest("GET", fmt.Sprintf("/%s", sId), nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		req = addContextHeaders(req)
		req = withURLParam(req, "id", sId)

		// ACT
		r.ServeHTTP(rr, req)

		//  ASSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

func TestAppserverChannelRolesHandler(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	sId := "123"
	cId := "456"

	r := chi.NewRouter()
	r.Get("/{sid}/channels/{cid}", api.AppserverChannelRolesHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	fullUrl := fmt.Sprintf("/%s/channels/%s", sId, cId)

	t.Run("Success:returns_list", func(t *testing.T) {
		mockService := new(testutil.MockChannelRoleService)
		mockRequest := &channel_role.ListChannelRolesRequest{
			ChannelId:   cId,
			AppserverId: sId,
		}
		mockResponse := &channel_role.ListChannelRolesResponse{
			ChannelRoles: []*channel_role.ChannelRole{
				{Id: "1", ChannelId: cId, AppserverId: sId, AppserverRoleId: "r1"},
			},
		}
		mockService.On("ListChannelRoles", mock.Anything, mockRequest).Return(mockResponse, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fullUrl, nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		expected := marshallResponse(t, api.CreateResponse([]types.ChannelRole{
			{ID: "1", ChannelId: cId, AppserverId: sId, AppserverRoleId: "r1"},
		}))
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		mockService := new(testutil.MockChannelRoleService)
		mockReq := &channel_role.ListChannelRolesRequest{
			ChannelId:   cId,
			AppserverId: sId,
		}
		mockService.On("ListChannelRoles", mock.Anything, mockReq).Return(nil, status.Error(codes.InvalidArgument, "invalid"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("GET", fullUrl, nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		expected := marshallResponse(t, api.CreateErrorResponse("invalid"))
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})
}

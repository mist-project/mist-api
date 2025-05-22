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
	pb_channel_role "mistapi/src/protos/v1/channel_role"
	"mistapi/src/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const channelRoleUrl = "/api/v1/channel-roles"

func TestCreateChannelRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:creates_channel_role", func(t *testing.T) {
		input := api.ChannelRoleCreate{
			ChannelId:       "c1",
			AppserverId:     "s1",
			AppserverRoleId: "r1",
		}
		mockReq := &pb_channel_role.CreateRequest{
			ChannelId:       input.ChannelId,
			AppserverId:     input.AppserverId,
			AppserverRoleId: input.AppserverRoleId,
		}
		mockResp := &pb_channel_role.CreateResponse{}

		mockService := new(testutil.MockChannelRoleService)
		mockService.On("Create", mock.Anything, mockReq).Return(mockResp, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		payload := marshallPayload(t, input)
		req, err := http.NewRequest("POST", channelRoleUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		api.ChannelRoleCreateHandler(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:invalid_payload", func(t *testing.T) {
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))

		req, err := http.NewRequest("POST", channelRoleUrl, marshallPayload(t, "invalid"))
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		api.ChannelRoleCreateHandler(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		input := api.ChannelRoleCreate{
			ChannelId:       "c2",
			AppserverId:     "s2",
			AppserverRoleId: "r2",
		}
		mockReq := &pb_channel_role.CreateRequest{
			ChannelId:       input.ChannelId,
			AppserverId:     input.AppserverId,
			AppserverRoleId: input.AppserverRoleId,
		}

		mockService := new(testutil.MockChannelRoleService)
		mockService.On("Create", mock.Anything, mockReq).Return(nil, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("POST", channelRoleUrl, marshallPayload(t, input))
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		api.ChannelRoleCreateHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListChannelRoles(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	channelId := "ch"
	appserverId := "srv"
	params := url.Values{}
	params.Add("channel_id", channelId)
	params.Add("appserver_id", appserverId)
	fullUrl := channelRoleUrl + "?" + params.Encode()

	r := chi.NewRouter()
	r.Get(channelRoleUrl, api.ChannelRoleListHandler)

	t.Run("Success:returns_list", func(t *testing.T) {
		mockService := new(testutil.MockChannelRoleService)
		mockRequest := &pb_channel_role.ListChannelRolesRequest{
			ChannelId:   channelId,
			AppserverId: appserverId,
		}
		mockResponse := &pb_channel_role.ListChannelRolesResponse{
			ChannelRoles: []*pb_channel_role.ChannelRole{
				{Id: "1", ChannelId: channelId, AppserverId: appserverId, AppserverRoleId: "r1"},
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

		expected := marshallResponse(t, api.CreateResponse([]api.ChannelRole{
			{ID: "1", ChannelId: channelId, AppserverId: appserverId, AppserverRoleId: "r1"},
		}))
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:missing_params", func(t *testing.T) {
		expected := marshallResponse(t, api.CreateErrorResponse("Channel ID and Appserver ID are required"))

		req, err := http.NewRequest("GET", channelRoleUrl, nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		mockService := new(testutil.MockChannelRoleService)
		mockReq := &pb_channel_role.ListChannelRolesRequest{
			ChannelId:   channelId,
			AppserverId: appserverId,
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

func TestDeleteChannelRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.ChannelRoleDeleteHandler)

	t.Run("Success:deletes_channel_role", func(t *testing.T) {
		id := "1"
		mockReq := &pb_channel_role.DeleteRequest{Id: id}
		mockResp := &pb_channel_role.DeleteResponse{}

		mockService := new(testutil.MockChannelRoleService)
		mockService.On("Delete", mock.Anything, mockReq).Return(mockResp, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", id), nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		req = withURLParam(req, "id", id)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		id := "fail"
		mockService := new(testutil.MockChannelRoleService)
		mockService.On("Delete", mock.Anything, &pb_channel_role.DeleteRequest{Id: id}).
			Return(nil, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", id), nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		req = withURLParam(req, "id", id)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

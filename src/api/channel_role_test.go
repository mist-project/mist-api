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
	"mistapi/src/protos/v1/channel_role"
	"mistapi/src/testutil"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const channelRoleUrl = "/api/v1/channel-roles"

func TestCreateChannelRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:creates_channel_role", func(t *testing.T) {
		// ARRANGE
		input := types.ChannelRoleCreate{
			ChannelId:       "c1",
			AppserverId:     "s1",
			AppserverRoleId: "r1",
		}
		mockReq := &channel_role.CreateRequest{
			ChannelId:       input.ChannelId,
			AppserverId:     input.AppserverId,
			AppserverRoleId: input.AppserverRoleId,
		}
		mockResp := &channel_role.CreateResponse{}

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

		// ACT
		api.ChannelRoleCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:invalid_payload", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))

		req, err := http.NewRequest("POST", channelRoleUrl, marshallPayload(t, "invalid"))
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.ChannelRoleCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		// ARRANGE
		input := types.ChannelRoleCreate{
			ChannelId:       "c2",
			AppserverId:     "s2",
			AppserverRoleId: "r2",
		}
		mockReq := &channel_role.CreateRequest{
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

		// ACT
		api.ChannelRoleCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteChannelRole(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.ChannelRoleDeleteHandler)

	t.Run("Success:deletes_channel_role", func(t *testing.T) {
		// ARRANGE
		id := "1"
		mockReq := &channel_role.DeleteRequest{Id: id}
		mockResp := &channel_role.DeleteResponse{}

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

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		// ARRANGE
		id := "fail"
		mockService := new(testutil.MockChannelRoleService)
		mockService.On("Delete", mock.Anything, &channel_role.DeleteRequest{Id: id}).
			Return(nil, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetChannelRoleClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", id), nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		req = withURLParam(req, "id", id)
		rr := httptest.NewRecorder()

		// ACT
		r.ServeHTTP(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

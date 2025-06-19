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
	"mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/testutil"
	"mistapi/src/types"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const roleSubUrl = "/api/v1/appserver-role-subs"

func TestCreateAppserverRoleSub(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:creates_role_sub", func(t *testing.T) {
		// ARRANGE
		input := types.AppserverRoleSubCreate{
			AppuserId:       "user",
			AppserverRoleId: "role",
			AppserverId:     "srv",
			AppserverSubId:  "sub",
		}
		mockReq := &appserver_role_sub.CreateRequest{
			AppuserId:       input.AppuserId,
			AppserverRoleId: input.AppserverRoleId,
			AppserverId:     input.AppserverId,
			AppserverSubId:  input.AppserverSubId,
		}
		mockResp := &appserver_role_sub.CreateResponse{}

		mockService := new(testutil.MockAppserverRoleSubService)
		mockService.On("Create", mock.Anything, mockReq).Return(mockResp, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		payload := marshallPayload(t, input)
		req, err := http.NewRequest("POST", roleSubUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverRoleSubCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:invalid_payload", func(t *testing.T) {
		// ARRANGE
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))

		payload := marshallPayload(t, "invalid")
		req, err := http.NewRequest("POST", roleSubUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverRoleSubCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_error", func(t *testing.T) {
		// ARRANGE
		input := types.AppserverRoleSubCreate{
			AppuserId:       "u",
			AppserverRoleId: "r",
			AppserverId:     "s",
			AppserverSubId:  "x",
		}
		mockReq := &appserver_role_sub.CreateRequest{
			AppuserId:       input.AppuserId,
			AppserverRoleId: input.AppserverRoleId,
			AppserverId:     input.AppserverId,
			AppserverSubId:  input.AppserverSubId,
		}
		mockService := new(testutil.MockAppserverRoleSubService)
		mockService.On("Create", mock.Anything, mockReq).Return(nil, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
		testutil.MockGrpcClient(t, mockClient)

		payload := marshallPayload(t, input)
		req, err := http.NewRequest("POST", roleSubUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		// ACT
		api.AppserverRoleSubCreateHandler(rr, req)

		// ASSERT
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteAppserverRoleSub(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.AppserverRoleSubDeleteHandler)

	t.Run("Success:deletes_role_sub", func(t *testing.T) {
		// ARRANGE
		id := "1"
		mockReq := &appserver_role_sub.DeleteRequest{Id: id}
		mockResp := &appserver_role_sub.DeleteResponse{}

		mockService := new(testutil.MockAppserverRoleSubService)
		mockService.On("Delete", mock.Anything, mockReq).Return(mockResp, nil)

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
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
		mockService := new(testutil.MockAppserverRoleSubService)
		mockService.On("Delete", mock.Anything, &appserver_role_sub.DeleteRequest{Id: id}).
			Return(nil, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
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

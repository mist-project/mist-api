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
	pb_appserver_role_sub "mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/testutil"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const roleSubUrl = "/api/v1/appserver-role-subs"

func TestCreateAppserverRoleSub(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	t.Run("Success:creates_role_sub", func(t *testing.T) {
		// ARRANGE
		input := api.AppserverRoleSubCreate{
			AppuserId:       "user",
			AppserverRoleId: "role",
			AppserverId:     "srv",
			AppserverSubId:  "sub",
		}
		mockReq := &pb_appserver_role_sub.CreateRequest{
			AppuserId:       input.AppuserId,
			AppserverRoleId: input.AppserverRoleId,
			AppserverId:     input.AppserverId,
			AppserverSubId:  input.AppserverSubId,
		}
		mockResp := &pb_appserver_role_sub.CreateResponse{}

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
		expected := marshallResponse(t, api.CreateErrorResponse("Invalid attributes provided."))

		payload := marshallPayload(t, "invalid")
		req, err := http.NewRequest("POST", roleSubUrl, payload)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		api.AppserverRoleSubCreateHandler(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_error", func(t *testing.T) {
		input := api.AppserverRoleSubCreate{
			AppuserId:       "u",
			AppserverRoleId: "r",
			AppserverId:     "s",
			AppserverSubId:  "x",
		}
		mockReq := &pb_appserver_role_sub.CreateRequest{
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

		api.AppserverRoleSubCreateHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListAppserverRoleSubs(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	appserverId := "srv"
	params := url.Values{}
	params.Add("appserver_id", appserverId)
	fullUrl := roleSubUrl + "?" + params.Encode()

	r := chi.NewRouter()
	r.Get(roleSubUrl, api.AppserverRoleSubListHandler)

	t.Run("Success:lists_role_subs", func(t *testing.T) {
		mockService := new(testutil.MockAppserverRoleSubService)
		mockRequest := &pb_appserver_role_sub.ListServerRoleSubsRequest{AppserverId: appserverId}
		mockResponse := &pb_appserver_role_sub.ListServerRoleSubsResponse{
			AppserverRoleSubs: []*pb_appserver_role_sub.AppserverRoleSub{
				{Id: "1", AppuserId: "user", AppserverRoleId: "role", AppserverId: appserverId},
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

		expected := marshallResponse(t, api.CreateResponse([]api.AppserverRoleSub{
			{ID: "1", AppuserId: "user", AppserverRoleId: "role", AppserverId: appserverId},
		}))

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:missing_appserver_id", func(t *testing.T) {
		expected := marshallResponse(t, api.CreateErrorResponse("Appserver ID is required"))

		req, err := http.NewRequest("GET", roleSubUrl, nil)
		require.NoError(t, err)
		req = addContextHeaders(req)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("Error:grpc_error", func(t *testing.T) {
		mockService := new(testutil.MockAppserverRoleSubService)
		mockRequest := &pb_appserver_role_sub.ListServerRoleSubsRequest{AppserverId: appserverId}
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

func TestDeleteAppserverRoleSub(t *testing.T) {
	log.SetOutput(new(strings.Builder))

	r := chi.NewRouter()
	r.Delete("/{id}", api.AppserverRoleSubDeleteHandler)

	t.Run("Success:deletes_role_sub", func(t *testing.T) {
		id := "1"
		mockReq := &pb_appserver_role_sub.DeleteRequest{Id: id}
		mockResp := &pb_appserver_role_sub.DeleteResponse{}

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

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("Error:grpc_failure", func(t *testing.T) {
		id := "fail"
		mockService := new(testutil.MockAppserverRoleSubService)
		mockService.On("Delete", mock.Anything, &pb_appserver_role_sub.DeleteRequest{Id: id}).
			Return(nil, errors.New("boom"))

		mockClient := new(testutil.MockClient)
		mockClient.On("GetAppserverRoleSubClient").Return(mockService)
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

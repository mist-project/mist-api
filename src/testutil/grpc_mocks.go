package testutil

import (
	"context"
	"testing"

	"mistapi/src/protos/v1/appserver"
	"mistapi/src/protos/v1/appserver_role"
	"mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/protos/v1/appserver_sub"
	"mistapi/src/protos/v1/appuser"
	"mistapi/src/protos/v1/channel"
	"mistapi/src/protos/v1/channel_role"

	"mistapi/src/service"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// returnIfError is our generic helper (assumed defined elsewhere in testutil).
// func returnIfError[T any](args mock.Arguments, index int) (T, error)

// ----- MESSAGES MOCKS -----
type MockClient struct {
	mock.Mock
}

func returnIfError[T any](args mock.Arguments, index int) (T, error) {
	err := args.Error(index)
	var zero T
	if err != nil {
		return zero, err
	}
	return args.Get(0).(T), nil
}

func (m *MockClient) GetAppserverClient() appserver.AppserverServiceClient {
	args := m.Called()
	return args.Get(0).(appserver.AppserverServiceClient)
}

func (m *MockClient) GetAppserverRoleClient() appserver_role.AppserverRoleServiceClient {
	args := m.Called()
	return args.Get(0).(appserver_role.AppserverRoleServiceClient)
}

func (m *MockClient) GetAppserverSubClient() appserver_sub.AppserverSubServiceClient {
	args := m.Called()
	return args.Get(0).(appserver_sub.AppserverSubServiceClient)
}

func (m *MockClient) GetAppserverRoleSubClient() appserver_role_sub.AppserverRoleSubServiceClient {
	args := m.Called()
	return args.Get(0).(appserver_role_sub.AppserverRoleSubServiceClient)
}

func (m *MockClient) GetAppuserClient() appuser.AppuserServiceClient {
	args := m.Called()
	return args.Get(0).(appuser.AppuserServiceClient)
}

func (m *MockClient) GetChannelClient() channel.ChannelServiceClient {
	args := m.Called()
	return args.Get(0).(channel.ChannelServiceClient)
}

func (m *MockClient) GetChannelRoleClient() channel_role.ChannelRoleServiceClient {
	args := m.Called()
	return args.Get(0).(channel_role.ChannelRoleServiceClient)
}

func MockGrpcClient(t *testing.T, mockClient service.GrpcClient) {
	original := service.NewGrpcClient
	service.NewGrpcClient = func() service.GrpcClient {
		return mockClient
	}
	t.Cleanup(func() {
		service.NewGrpcClient = original
	})
}

// ----- GRPC MOCKS -----
type MockAppserverService struct{ mock.Mock }

func (m *MockAppserverService) Create(
	ctx context.Context, in *appserver.CreateRequest, opts ...grpc.CallOption,
) (*appserver.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver.CreateResponse](args, 1)
}

func (m *MockAppserverService) GetById(
	ctx context.Context, in *appserver.GetByIdRequest, opts ...grpc.CallOption,
) (*appserver.GetByIdResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver.GetByIdResponse](args, 1)
}

func (m *MockAppserverService) List(
	ctx context.Context, in *appserver.ListRequest, opts ...grpc.CallOption,
) (*appserver.ListResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver.ListResponse](args, 1)
}

func (m *MockAppserverService) Delete(
	ctx context.Context, in *appserver.DeleteRequest, opts ...grpc.CallOption,
) (*appserver.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver.DeleteResponse](args, 1)
}

type MockAppserverPermissionService struct{ mock.Mock }

// ----- APPSERVER ROLE -----
type MockAppserverRoleService struct{ mock.Mock }

func (m *MockAppserverRoleService) Create(
	ctx context.Context, in *appserver_role.CreateRequest, opts ...grpc.CallOption,
) (*appserver_role.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_role.CreateResponse](args, 1)
}

func (m *MockAppserverRoleService) ListServerRoles(
	ctx context.Context, in *appserver_role.ListServerRolesRequest, opts ...grpc.CallOption,
) (*appserver_role.ListServerRolesResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_role.ListServerRolesResponse](args, 1)
}

func (m *MockAppserverRoleService) Delete(
	ctx context.Context, in *appserver_role.DeleteRequest, opts ...grpc.CallOption,
) (*appserver_role.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_role.DeleteResponse](args, 1)
}

// ----- APPSERVER ROLE SUB -----
type MockAppserverRoleSubService struct{ mock.Mock }

func (m *MockAppserverRoleSubService) Create(
	ctx context.Context, in *appserver_role_sub.CreateRequest, opts ...grpc.CallOption,
) (*appserver_role_sub.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_role_sub.CreateResponse](args, 1)
}

func (m *MockAppserverRoleSubService) ListServerRoleSubs(
	ctx context.Context, in *appserver_role_sub.ListServerRoleSubsRequest, opts ...grpc.CallOption,
) (*appserver_role_sub.ListServerRoleSubsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_role_sub.ListServerRoleSubsResponse](args, 1)
}

func (m *MockAppserverRoleSubService) Delete(
	ctx context.Context, in *appserver_role_sub.DeleteRequest, opts ...grpc.CallOption,
) (*appserver_role_sub.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_role_sub.DeleteResponse](args, 1)
}

// ----- APPSERVER SUB -----
type MockAppserverSubService struct{ mock.Mock }

func (m *MockAppserverSubService) Create(
	ctx context.Context, in *appserver_sub.CreateRequest, opts ...grpc.CallOption,
) (*appserver_sub.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_sub.CreateResponse](args, 1)
}

func (m *MockAppserverSubService) ListUserServerSubs(
	ctx context.Context, in *appserver_sub.ListUserServerSubsRequest, opts ...grpc.CallOption,
) (*appserver_sub.ListUserServerSubsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_sub.ListUserServerSubsResponse](args, 1)
}

func (m *MockAppserverSubService) ListAppserverUserSubs(
	ctx context.Context, in *appserver_sub.ListAppserverUserSubsRequest, opts ...grpc.CallOption,
) (*appserver_sub.ListAppserverUserSubsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_sub.ListAppserverUserSubsResponse](args, 1)
}

func (m *MockAppserverSubService) Delete(
	ctx context.Context, in *appserver_sub.DeleteRequest, opts ...grpc.CallOption,
) (*appserver_sub.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*appserver_sub.DeleteResponse](args, 1)
}

// ----- CHANNEL -----
type MockChannelService struct{ mock.Mock }

func (m *MockChannelService) Create(
	ctx context.Context, in *channel.CreateRequest, opts ...grpc.CallOption,
) (*channel.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel.CreateResponse](args, 1)
}

func (m *MockChannelService) GetById(
	ctx context.Context, in *channel.GetByIdRequest, opts ...grpc.CallOption,
) (*channel.GetByIdResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel.GetByIdResponse](args, 1)
}

func (m *MockChannelService) ListServerChannels(
	ctx context.Context, in *channel.ListServerChannelsRequest, opts ...grpc.CallOption,
) (*channel.ListServerChannelsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel.ListServerChannelsResponse](args, 1)
}

func (m *MockChannelService) Delete(
	ctx context.Context, in *channel.DeleteRequest, opts ...grpc.CallOption,
) (*channel.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel.DeleteResponse](args, 1)
}

// ----- CHANNEL ROLE -----
type MockChannelRoleService struct{ mock.Mock }

func (m *MockChannelRoleService) Create(
	ctx context.Context, in *channel_role.CreateRequest, opts ...grpc.CallOption,
) (*channel_role.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel_role.CreateResponse](args, 1)
}

func (m *MockChannelRoleService) ListChannelRoles(
	ctx context.Context, in *channel_role.ListChannelRolesRequest, opts ...grpc.CallOption,
) (*channel_role.ListChannelRolesResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel_role.ListChannelRolesResponse](args, 1)
}

func (m *MockChannelRoleService) Delete(
	ctx context.Context, in *channel_role.DeleteRequest, opts ...grpc.CallOption,
) (*channel_role.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*channel_role.DeleteResponse](args, 1)
}

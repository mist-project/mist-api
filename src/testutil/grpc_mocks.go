package testutil

import (
	"context"
	"testing"

	pb_appserver "mistapi/src/protos/v1/appserver"
	pb_appserverpermission "mistapi/src/protos/v1/appserver_permission"
	pb_appserver_role "mistapi/src/protos/v1/appserver_role"
	pb_appserver_rolesub "mistapi/src/protos/v1/appserver_role_sub"
	pb_appserver_sub "mistapi/src/protos/v1/appserver_sub"
	pb_appuser "mistapi/src/protos/v1/appuser"
	pb_channel "mistapi/src/protos/v1/channel"
	pb_channelrole "mistapi/src/protos/v1/channel_role"

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

func (m *MockClient) GetAppserverClient() pb_appserver.AppserverServiceClient {
	args := m.Called()
	return args.Get(0).(pb_appserver.AppserverServiceClient)
}

func (m *MockClient) GetAppserverRoleClient() pb_appserver_role.AppserverRoleServiceClient {
	args := m.Called()
	return args.Get(0).(pb_appserver_role.AppserverRoleServiceClient)
}

func (m *MockClient) GetAppserverSubClient() pb_appserver_sub.AppserverSubServiceClient {
	args := m.Called()
	return args.Get(0).(pb_appserver_sub.AppserverSubServiceClient)
}

func (m *MockClient) GetAppserverRoleSubClient() pb_appserver_rolesub.AppserverRoleSubServiceClient {
	args := m.Called()
	return args.Get(0).(pb_appserver_rolesub.AppserverRoleSubServiceClient)
}

func (m *MockClient) GetAppserverPermissionClient() pb_appserverpermission.AppserverPermissionServiceClient {
	args := m.Called()
	return args.Get(0).(pb_appserverpermission.AppserverPermissionServiceClient)
}

func (m *MockClient) GetAppuserClient() pb_appuser.AppuserServiceClient {
	args := m.Called()
	return args.Get(0).(pb_appuser.AppuserServiceClient)
}

func (m *MockClient) GetChannelClient() pb_channel.ChannelServiceClient {
	args := m.Called()
	return args.Get(0).(pb_channel.ChannelServiceClient)
}

func (m *MockClient) GetChannelRoleClient() pb_channelrole.ChannelRoleServiceClient {
	args := m.Called()
	return args.Get(0).(pb_channelrole.ChannelRoleServiceClient)
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
	ctx context.Context, in *pb_appserver.CreateRequest, opts ...grpc.CallOption,
) (*pb_appserver.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver.CreateResponse](args, 1)
}

func (m *MockAppserverService) GetById(
	ctx context.Context, in *pb_appserver.GetByIdRequest, opts ...grpc.CallOption,
) (*pb_appserver.GetByIdResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver.GetByIdResponse](args, 1)
}

func (m *MockAppserverService) List(
	ctx context.Context, in *pb_appserver.ListRequest, opts ...grpc.CallOption,
) (*pb_appserver.ListResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver.ListResponse](args, 1)
}

func (m *MockAppserverService) Delete(
	ctx context.Context, in *pb_appserver.DeleteRequest, opts ...grpc.CallOption,
) (*pb_appserver.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver.DeleteResponse](args, 1)
}

type MockAppserverPermissionService struct{ mock.Mock }

func (m *MockAppserverPermissionService) Create(
	ctx context.Context, in *pb_appserverpermission.CreateRequest, opts ...grpc.CallOption,
) (*pb_appserverpermission.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserverpermission.CreateResponse](args, 1)
}

func (m *MockAppserverPermissionService) ListAppserverUsers(
	ctx context.Context, in *pb_appserverpermission.ListAppserverUsersRequest, opts ...grpc.CallOption,
) (*pb_appserverpermission.ListAppserverUsersResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserverpermission.ListAppserverUsersResponse](args, 1)
}

func (m *MockAppserverPermissionService) Delete(
	ctx context.Context, in *pb_appserverpermission.DeleteRequest, opts ...grpc.CallOption,
) (*pb_appserverpermission.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserverpermission.DeleteResponse](args, 1)
}

// ----- APPSERVER ROLE -----
type MockAppserverRoleService struct{ mock.Mock }

func (m *MockAppserverRoleService) Create(
	ctx context.Context, in *pb_appserver_role.CreateRequest, opts ...grpc.CallOption,
) (*pb_appserver_role.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_role.CreateResponse](args, 1)
}

func (m *MockAppserverRoleService) ListServerRoles(
	ctx context.Context, in *pb_appserver_role.ListServerRolesRequest, opts ...grpc.CallOption,
) (*pb_appserver_role.ListServerRolesResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_role.ListServerRolesResponse](args, 1)
}

func (m *MockAppserverRoleService) Delete(
	ctx context.Context, in *pb_appserver_role.DeleteRequest, opts ...grpc.CallOption,
) (*pb_appserver_role.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_role.DeleteResponse](args, 1)
}

// ----- APPSERVER ROLE SUB -----
type MockAppserverRoleSubService struct{ mock.Mock }

func (m *MockAppserverRoleSubService) Create(
	ctx context.Context, in *pb_appserver_rolesub.CreateRequest, opts ...grpc.CallOption,
) (*pb_appserver_rolesub.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_rolesub.CreateResponse](args, 1)
}

func (m *MockAppserverRoleSubService) ListServerRoleSubs(
	ctx context.Context, in *pb_appserver_rolesub.ListServerRoleSubsRequest, opts ...grpc.CallOption,
) (*pb_appserver_rolesub.ListServerRoleSubsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_rolesub.ListServerRoleSubsResponse](args, 1)
}

func (m *MockAppserverRoleSubService) Delete(
	ctx context.Context, in *pb_appserver_rolesub.DeleteRequest, opts ...grpc.CallOption,
) (*pb_appserver_rolesub.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_rolesub.DeleteResponse](args, 1)
}

// ----- APPSERVER SUB -----
type MockAppserverSubService struct{ mock.Mock }

func (m *MockAppserverSubService) Create(
	ctx context.Context, in *pb_appserver_sub.CreateRequest, opts ...grpc.CallOption,
) (*pb_appserver_sub.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_sub.CreateResponse](args, 1)
}

func (m *MockAppserverSubService) ListUserServerSubs(
	ctx context.Context, in *pb_appserver_sub.ListUserServerSubsRequest, opts ...grpc.CallOption,
) (*pb_appserver_sub.ListUserServerSubsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_sub.ListUserServerSubsResponse](args, 1)
}

func (m *MockAppserverSubService) ListAppserverUserSubs(
	ctx context.Context, in *pb_appserver_sub.ListAppserverUserSubsRequest, opts ...grpc.CallOption,
) (*pb_appserver_sub.ListAppserverUserSubsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_sub.ListAppserverUserSubsResponse](args, 1)
}

func (m *MockAppserverSubService) Delete(
	ctx context.Context, in *pb_appserver_sub.DeleteRequest, opts ...grpc.CallOption,
) (*pb_appserver_sub.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_appserver_sub.DeleteResponse](args, 1)
}

// ----- CHANNEL -----
type MockChannelService struct{ mock.Mock }

func (m *MockChannelService) Create(
	ctx context.Context, in *pb_channel.CreateRequest, opts ...grpc.CallOption,
) (*pb_channel.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channel.CreateResponse](args, 1)
}

func (m *MockChannelService) GetById(
	ctx context.Context, in *pb_channel.GetByIdRequest, opts ...grpc.CallOption,
) (*pb_channel.GetByIdResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channel.GetByIdResponse](args, 1)
}

func (m *MockChannelService) ListServerChannels(
	ctx context.Context, in *pb_channel.ListServerChannelsRequest, opts ...grpc.CallOption,
) (*pb_channel.ListServerChannelsResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channel.ListServerChannelsResponse](args, 1)
}

func (m *MockChannelService) Delete(
	ctx context.Context, in *pb_channel.DeleteRequest, opts ...grpc.CallOption,
) (*pb_channel.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channel.DeleteResponse](args, 1)
}

// ----- CHANNEL ROLE -----
type MockChannelRoleService struct{ mock.Mock }

func (m *MockChannelRoleService) Create(
	ctx context.Context, in *pb_channelrole.CreateRequest, opts ...grpc.CallOption,
) (*pb_channelrole.CreateResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channelrole.CreateResponse](args, 1)
}

func (m *MockChannelRoleService) ListChannelRoles(
	ctx context.Context, in *pb_channelrole.ListChannelRolesRequest, opts ...grpc.CallOption,
) (*pb_channelrole.ListChannelRolesResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channelrole.ListChannelRolesResponse](args, 1)
}

func (m *MockChannelRoleService) Delete(
	ctx context.Context, in *pb_channelrole.DeleteRequest, opts ...grpc.CallOption,
) (*pb_channelrole.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return returnIfError[*pb_channelrole.DeleteResponse](args, 1)
}

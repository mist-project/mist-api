package api_test

import (
	"context"
	"testing"

	pb "mistapi/src/protos/v1/gen"
	"mistapi/src/service"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockClient struct {
	mock.Mock
}

type MockService struct{ mock.Mock }

// ----- MESSAGES MOCKS -----
func (m *MockClient) GetServerClient() pb.AppserverServiceClient {
	args := m.Called()
	return args.Get(0).(pb.AppserverServiceClient)
}
func (m *MockClient) GetChannelClient() pb.ChannelServiceClient {
	args := m.Called()
	return args.Get(0).(pb.ChannelServiceClient)
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

// ----- GRPC MOCKS ----
func (m *MockService) CreateAppserver(ctx context.Context, in *pb.CreateAppserverRequest, opts ...grpc.CallOption,
) (*pb.CreateAppserverResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CreateAppserverResponse), args.Error(1)
}

func (m *MockService) GetByIdAppserver(ctx context.Context, in *pb.GetByIdAppserverRequest, opts ...grpc.CallOption,
) (*pb.GetByIdAppserverResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GetByIdAppserverResponse), args.Error(1)
}
func (m *MockService) ListAppservers(ctx context.Context, in *pb.ListAppserversRequest, opts ...grpc.CallOption,
) (*pb.ListAppserversResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ListAppserversResponse), args.Error(1)
}
func (m *MockService) DeleteAppserver(ctx context.Context, in *pb.DeleteAppserverRequest, opts ...grpc.CallOption,
) (*pb.DeleteAppserverResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DeleteAppserverResponse), args.Error(1)
}

// ----- APPSERVER SUB -----
func (m *MockService) CreateAppserverSub(ctx context.Context, in *pb.CreateAppserverSubRequest, opts ...grpc.CallOption,
) (*pb.CreateAppserverSubResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CreateAppserverSubResponse), args.Error(1)
}
func (m *MockService) GetUserAppserverSubs(ctx context.Context, in *pb.GetUserAppserverSubsRequest, opts ...grpc.CallOption,
) (*pb.GetUserAppserverSubsResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GetUserAppserverSubsResponse), args.Error(1)
}
func (m *MockService) DeleteAppserverSub(ctx context.Context, in *pb.DeleteAppserverSubRequest, opts ...grpc.CallOption,
) (*pb.DeleteAppserverSubResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DeleteAppserverSubResponse), args.Error(1)
}
func (m *MockService) GetAllUsersAppserverSubs(ctx context.Context, in *pb.GetAllUsersAppserverSubsRequest, opts ...grpc.CallOption,
) (*pb.GetAllUsersAppserverSubsResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GetAllUsersAppserverSubsResponse), args.Error(1)
}

// ----- APPSERVER ROLE -----
func (m *MockService) CreateAppserverRole(ctx context.Context, in *pb.CreateAppserverRoleRequest, opts ...grpc.CallOption,
) (*pb.CreateAppserverRoleResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CreateAppserverRoleResponse), args.Error(1)
}
func (m *MockService) GetAllAppserverRoles(ctx context.Context, in *pb.GetAllAppserverRolesRequest, opts ...grpc.CallOption,
) (*pb.GetAllAppserverRolesResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GetAllAppserverRolesResponse), args.Error(1)
}
func (m *MockService) DeleteAppserverRole(ctx context.Context, in *pb.DeleteAppserverRoleRequest, opts ...grpc.CallOption,
) (*pb.DeleteAppserverRoleResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DeleteAppserverRoleResponse), args.Error(1)
}

// ----- APPSERVER ROLE SUB -----
func (m *MockService) CreateAppserverRoleSub(ctx context.Context, in *pb.CreateAppserverRoleSubRequest, opts ...grpc.CallOption,
) (*pb.CreateAppserverRoleSubResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CreateAppserverRoleSubResponse), args.Error(1)
}
func (m *MockService) DeleteAppserverRoleSub(ctx context.Context, in *pb.DeleteAppserverRoleSubRequest, opts ...grpc.CallOption,
) (*pb.DeleteAppserverRoleSubResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DeleteAppserverRoleSubResponse), args.Error(1)
}

// ----- CHANNEL -----
func (m *MockService) CreateChannel(ctx context.Context, in *pb.CreateChannelRequest, opts ...grpc.CallOption,
) (*pb.CreateChannelResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CreateChannelResponse), args.Error(1)
}
func (m *MockService) ListChannels(ctx context.Context, in *pb.ListChannelsRequest, opts ...grpc.CallOption,
) (*pb.ListChannelsResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ListChannelsResponse), args.Error(1)
}
func (m *MockService) DeleteChannel(ctx context.Context, in *pb.DeleteChannelRequest, opts ...grpc.CallOption,
) (*pb.DeleteChannelResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DeleteChannelResponse), args.Error(1)
}
func (m *MockService) GetByIdChannel(ctx context.Context, in *pb.GetByIdChannelRequest, opts ...grpc.CallOption,
) (*pb.GetByIdChannelResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GetByIdChannelResponse), args.Error(1)
}

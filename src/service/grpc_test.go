package service_test

import (
	"fmt"
	"testing"

	"mistapi/src/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetAppserverClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverClient := client.GetAppserverClient()

	// ASSERT
	assert.NotNil(t, serverClient)
}

func TestGetAppserverRoleClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverRoleClient := client.GetAppserverRoleClient()

	// ASSERT
	assert.NotNil(t, serverRoleClient)
}

func TestGetAppserverSubClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverSubClient := client.GetAppserverSubClient()

	// ASSERT
	assert.NotNil(t, serverSubClient)
}

func TestGetAppserverRoleSubClient(t *testing.T) {

	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverRoleSubClient := client.GetAppserverRoleSubClient()

	// ASSERT
	assert.NotNil(t, serverRoleSubClient)
}

func TestGetAppserverPermissionClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverPermissionClient := client.GetAppserverPermissionClient()

	// ASSERT
	assert.NotNil(t, serverPermissionClient)
}

func TestGetAppuserClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	appUserClient := client.GetAppuserClient()

	// ASSERT
	assert.NotNil(t, appUserClient)
}

func TestGetChannelClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverClient := client.GetChannelClient()

	// ASSERT
	assert.NotNil(t, serverClient)
}

func TestGetChannelRoleClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	channelRoleClient := client.GetChannelRoleClient()

	// ASSERT
	assert.NotNil(t, channelRoleClient)
}

func TestSetupGrpcHeaders(t *testing.T) {
	jwt := "test-jwt-token"
	ctx, cancel := service.SetupGrpcHeaders(jwt)
	defer cancel()

	md, ok := metadata.FromOutgoingContext(ctx)
	require.True(t, ok, "expected metadata to be in outgoing context")

	authHeader := md.Get("authorization")
	require.Len(t, authHeader, 1, "expected one authorization header")
	assert.Equal(t, fmt.Sprintf("Bearer %s", jwt), authHeader[0])
}

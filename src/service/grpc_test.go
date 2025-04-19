package service_test

import (
	"fmt"
	"mistapi/src/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetServerClient(t *testing.T) {
	// ARRANGE
	mockConn := new(grpc.ClientConn)
	client := service.Client{
		Conn: mockConn,
	}

	// ACT
	serverClient := client.GetServerClient()

	// ASSERT
	assert.NotNil(t, serverClient)
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

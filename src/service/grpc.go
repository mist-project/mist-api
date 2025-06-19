package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"mistapi/src/protos/v1/appserver"
	"mistapi/src/protos/v1/appserver_role"
	"mistapi/src/protos/v1/appserver_role_sub"
	"mistapi/src/protos/v1/appserver_sub"
	"mistapi/src/protos/v1/appuser"
	"mistapi/src/protos/v1/channel"
	"mistapi/src/protos/v1/channel_role"
)

type GrpcClient interface {
	GetAppserverClient() appserver.AppserverServiceClient
	GetAppserverRoleClient() appserver_role.AppserverRoleServiceClient
	GetAppserverSubClient() appserver_sub.AppserverSubServiceClient
	GetAppserverRoleSubClient() appserver_role_sub.AppserverRoleSubServiceClient
	GetAppuserClient() appuser.AppuserServiceClient
	GetChannelClient() channel.ChannelServiceClient
	GetChannelRoleClient() channel_role.ChannelRoleServiceClient
}

type Client struct {
	Conn *grpc.ClientConn
}

var (
	conn     *grpc.ClientConn
	connOnce sync.Once
)

func (c Client) GetAppserverClient() appserver.AppserverServiceClient {
	return appserver.NewAppserverServiceClient(c.Conn)
}

func (c Client) GetAppserverRoleClient() appserver_role.AppserverRoleServiceClient {
	return appserver_role.NewAppserverRoleServiceClient(c.Conn)
}

func (c Client) GetAppserverSubClient() appserver_sub.AppserverSubServiceClient {
	return appserver_sub.NewAppserverSubServiceClient(c.Conn)
}

func (c Client) GetAppserverRoleSubClient() appserver_role_sub.AppserverRoleSubServiceClient {
	return appserver_role_sub.NewAppserverRoleSubServiceClient(c.Conn)
}

func (c Client) GetAppuserClient() appuser.AppuserServiceClient {
	return appuser.NewAppuserServiceClient(c.Conn)
}

func (c Client) GetChannelClient() channel.ChannelServiceClient {
	return channel.NewChannelServiceClient(c.Conn)
}

func (c Client) GetChannelRoleClient() channel_role.ChannelRoleServiceClient {
	return channel_role.NewChannelRoleServiceClient(c.Conn)
}

func SetupGrpcHeaders(jwtT string) (context.Context, context.CancelFunc) {
	// TODO: replace timeout in the future
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	grpcMetadata := metadata.Pairs(
		"authorization", fmt.Sprintf("Bearer %s", jwtT),
	)

	ctx = metadata.NewOutgoingContext(ctx, grpcMetadata)
	return ctx, cancel
}

func GetGrpcClientConnection() *grpc.ClientConn {
	connOnce.Do(func() {
		var err error
		conn, err = grpc.NewClient(
			os.Getenv("MIST_BACKEND_APP_URL"),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Panicf("Error communicating with backend service: %v", err)
		}
	})
	return conn
}

var NewGrpcClient = func() GrpcClient {
	return Client{Conn: GetGrpcClientConnection()}
}

func CloseGrpcConnection() {
	if conn != nil {
		_ = conn.Close()
	}
}

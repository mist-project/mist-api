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

	pb_appserver "mistapi/src/protos/v1/appserver"
	pb_appserverpermission "mistapi/src/protos/v1/appserver_permission"
	pb_appserver_role "mistapi/src/protos/v1/appserver_role"
	pb_appserver_rolesub "mistapi/src/protos/v1/appserver_role_sub"
	pb_appserver_sub "mistapi/src/protos/v1/appserver_sub"
	pb_appuser "mistapi/src/protos/v1/appuser"
	pb_channel "mistapi/src/protos/v1/channel"
	pb_channelrole "mistapi/src/protos/v1/channel_role"
)

type GrpcClient interface {
	GetAppserverClient() pb_appserver.AppserverServiceClient
	GetAppserverRoleClient() pb_appserver_role.AppserverRoleServiceClient
	GetAppserverSubClient() pb_appserver_sub.AppserverSubServiceClient
	GetAppserverRoleSubClient() pb_appserver_rolesub.AppserverRoleSubServiceClient
	GetAppserverPermissionClient() pb_appserverpermission.AppserverPermissionServiceClient
	GetAppuserClient() pb_appuser.AppuserServiceClient
	GetChannelClient() pb_channel.ChannelServiceClient
	GetChannelRoleClient() pb_channelrole.ChannelRoleServiceClient
}

type Client struct {
	Conn *grpc.ClientConn
}

var (
	conn     *grpc.ClientConn
	connOnce sync.Once
)

func (c Client) GetAppserverClient() pb_appserver.AppserverServiceClient {
	return pb_appserver.NewAppserverServiceClient(c.Conn)
}

func (c Client) GetAppserverRoleClient() pb_appserver_role.AppserverRoleServiceClient {
	return pb_appserver_role.NewAppserverRoleServiceClient(c.Conn)
}

func (c Client) GetAppserverSubClient() pb_appserver_sub.AppserverSubServiceClient {
	return pb_appserver_sub.NewAppserverSubServiceClient(c.Conn)
}

func (c Client) GetAppserverRoleSubClient() pb_appserver_rolesub.AppserverRoleSubServiceClient {
	return pb_appserver_rolesub.NewAppserverRoleSubServiceClient(c.Conn)
}

func (c Client) GetAppserverPermissionClient() pb_appserverpermission.AppserverPermissionServiceClient {
	return pb_appserverpermission.NewAppserverPermissionServiceClient(c.Conn)
}

func (c Client) GetAppuserClient() pb_appuser.AppuserServiceClient {
	return pb_appuser.NewAppuserServiceClient(c.Conn)
}

func (c Client) GetChannelClient() pb_channel.ChannelServiceClient {
	return pb_channel.NewChannelServiceClient(c.Conn)
}

func (c Client) GetChannelRoleClient() pb_channelrole.ChannelRoleServiceClient {
	return pb_channelrole.NewChannelRoleServiceClient(c.Conn)
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

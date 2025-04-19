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

	pb "mistapi/src/protos/v1/gen"
)

type GrpcClient interface {
	GetServerClient() pb.AppserverServiceClient
	GetChannelClient() pb.ChannelServiceClient
}

type Client struct {
	Conn *grpc.ClientConn
}

var (
	conn     *grpc.ClientConn
	connOnce sync.Once
)

func (c Client) GetServerClient() pb.AppserverServiceClient {
	return pb.NewAppserverServiceClient(c.Conn)
}

func (c Client) GetChannelClient() pb.ChannelServiceClient {
	return pb.NewChannelServiceClient(c.Conn)
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

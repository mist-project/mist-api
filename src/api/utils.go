package api

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
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

func setupContext(jwtT string) (context.Context, context.CancelFunc) {
	// TODO: replace timeout in the future
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	grpcMetadata := metadata.Pairs(
		"authorization", fmt.Sprintf("Bearer %s", jwtT),
	)

	ctx = metadata.NewOutgoingContext(ctx, grpcMetadata)
	return ctx, cancel
}

func (c Client) GetServerClient() pb.AppserverServiceClient {
	return pb.NewAppserverServiceClient(c.Conn)
}

func (c Client) GetChannelClient() pb.ChannelServiceClient {
	return pb.NewChannelServiceClient(c.Conn)
}

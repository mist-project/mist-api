package api

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "mistapi/src/protos/v1/gen"
)

var tempT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNzQ0ODUzMDExLCJpYXQiOjE3NDQ4NTEyMTEsImp0aSI6IjBjOWJkMDgyMDI5NzRiZDViODFhZWRhZTk2NjFlNDQ4IiwidXNlcl9pZCI6IjdkOWI1YjVhLTg5YjMtNDQxNC04YzE5LWY3YzJkNTQ0YjBjYSIsImF1ZCI6WyJtaXN0LWJhY2tlbmQiLCJtaXN0LWFwaSIsIm1pc3QtaW8iXSwiaXNzIjoibWlzdC1hcGkifQ.tSWc8npSUcYSn6kwCNq_eedOAy58ZWcdZqWIK53JRXc"

func setupContext() (context.Context, context.CancelFunc) {
	// TODO: replace timeout in the future
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	grpcMetadata := metadata.Pairs(
		"authorization", fmt.Sprintf("Bearer %s", tempT),
	)

	ctx = metadata.NewOutgoingContext(ctx, grpcMetadata)
	return ctx, cancel
}

type GrpcClient interface {
	GetServerClient() pb.AppserverServiceClient
	GetChannelClient() pb.ChannelServiceClient
}

type Client struct {
	Conn *grpc.ClientConn
}

func (c Client) GetServerClient() pb.AppserverServiceClient {
	return pb.NewAppserverServiceClient(c.Conn)
}

func (c Client) GetChannelClient() pb.ChannelServiceClient {
	return pb.NewChannelServiceClient(c.Conn)
}

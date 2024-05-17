package sub

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/sub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hrvadl/btcratenotifier/gw/pkg/logger"
)

func NewClient(addr string, log *slog.Logger) (*Client, error) {
	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logger.NewClientGRPCMiddleware(log),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sender service: %w", err)
	}

	return &Client{
		api: pb.NewSubServiceClient(cc),
	}, nil
}

type Client struct {
	api pb.SubServiceClient
}

func (c *Client) Subscribe(ctx context.Context, req *pb.SubscribeRequest) error {
	_, err := c.api.Subscribe(ctx, req)
	return err
}

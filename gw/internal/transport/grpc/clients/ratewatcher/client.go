package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

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
		api: pb.NewRateWatcherServiceClient(cc),
		log: log,
	}, nil
}

type Client struct {
	log *slog.Logger
	api pb.RateWatcherServiceClient
}

func (c *Client) GetRate(ctx context.Context) (float32, error) {
	resp, err := c.api.GetRate(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}

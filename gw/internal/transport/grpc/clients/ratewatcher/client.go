package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang/protobuf/protoc-gen-go/grpc"
	pb "github.com/hrvadl/protos/gen/go/v1/ratewatcher"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(ctx context.Context, addr string, log *slog.Logger) (*Client, error) {
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	resp, err := c.api.GetRate(ctx)
	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}

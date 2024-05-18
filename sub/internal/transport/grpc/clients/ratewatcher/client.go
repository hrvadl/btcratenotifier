package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hrvadl/btcratenotifier/sub/pkg/logger"
)

const (
	retryCount   = 3
	retryTimeout = time.Second * 2
)

func NewClient(addr string, log *slog.Logger) (*Client, error) {
	retryOpt := []retry.CallOption{
		retry.WithCodes(codes.Aborted, codes.NotFound, codes.DeadlineExceeded),
		retry.WithMax(retryCount),
		retry.WithPerRetryTimeout(retryTimeout),
	}

	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logger.NewClientGRPCMiddleware(log),
			retry.UnaryClientInterceptor(retryOpt...),
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

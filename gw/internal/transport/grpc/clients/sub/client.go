//go:generate mockgen -destination=./mocks/mock_sub.go -package=mocks github.com/hrvadl/converter/protos/gen/go/v1/sub SubServiceClient
package sub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	pb "github.com/hrvadl/converter/protos/gen/go/v1/sub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hrvadl/converter/gw/pkg/logger"
)

const (
	retryCount   = 3
	retryTimeout = time.Second * 2
)

// NewClient constructs a GRPC subscriber client with provided arguments. Under the hood
// it initializes a bunch of GRPC middleware for debugging and monitoring purposes. I.E:
// - retry middleware
// - request logger middleware
// If initialization of connection has failed it will return an error.
// NOTE: neither of parameters couldn't be nil or client will panic.
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
		api: pb.NewSubServiceClient(cc),
	}, nil
}

// Client represents GRPC subscriber client which
// is responsible for subscribing new users and triggering
// email notification to the subscribers.
type Client struct {
	api pb.SubServiceClient
}

func (c *Client) Subscribe(ctx context.Context, req *pb.SubscribeRequest) error {
	_, err := c.api.Subscribe(ctx, req)
	return err
}

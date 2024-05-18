//go:generate mockgen -destination=./mocks/mock_parser.go -package=mocks github.com/hrvadl/converter/protos/gen/go/v1/mailer MailerServiceClient
package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	pb "github.com/hrvadl/converter/protos/gen/go/v1/mailer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hrvadl/converter/sub/pkg/logger"
)

const (
	retryCount   = 3
	retryTimeout = time.Second * 2
)

func NewClient(addr string, from string, log *slog.Logger) (*Client, error) {
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
		api:  pb.NewMailerServiceClient(cc),
		log:  log,
		from: from,
	}, nil
}

type Client struct {
	log  *slog.Logger
	api  pb.MailerServiceClient
	from string
}

type SendOptions struct {
	To      string
	Payload string
}

func (c *Client) Send(ctx context.Context, html, subject string, to ...string) error {
	_, err := c.api.Send(ctx, &pb.Mail{
		From:    c.from,
		To:      to,
		Subject: subject,
		Html:    html,
	})
	return err
}

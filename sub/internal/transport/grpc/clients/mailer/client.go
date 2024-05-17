package mailer

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/mailer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hrvadl/btcratenotifier/sub/pkg/logger"
)

func NewClient(addr string, from string, log *slog.Logger) (*Client, error) {
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

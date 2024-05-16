package mailer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hrvadl/btcratenotifier/pkg/logger"
	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/mailer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		api: pb.NewMailerServiceClient(cc),
		log: log,
	}, nil
}

type Client struct {
	log *slog.Logger
	api pb.MailerServiceClient
}

type SendOptions struct {
	To      string
	Payload string
}

func (c *Client) Send(ctx context.Context, html string, to ...string) error {
	_, err := c.api.Send(ctx, &pb.Mail{
		From:    "vadym@test.com",
		To:      to,
		Subject: "BTC to UAH rate exchange",
		Html:    html,
	})
	return err
}

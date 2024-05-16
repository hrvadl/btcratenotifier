package mailer

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/mailer"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(addr string, log *slog.Logger) (*Client, error) {
	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (c *Client) Send(ctx context.Context, emails ...SendOptions) error {
	stream, err := c.api.Send(ctx)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	for _, m := range emails {
		g.Go(func() error {
			return stream.Send(&pb.Mail{
				Recipient: m.To,
				Payload:   m.Payload,
			})
		})
	}

	return g.Wait()
}

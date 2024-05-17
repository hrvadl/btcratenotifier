package resend

import (
	"context"
	"fmt"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/mailer"
	rs "github.com/resend/resend-go/v2"
)

const operation = "resend mail client"

func NewClient(token string) *Client {
	return &Client{
		c: rs.NewClient(token),
	}
}

type Client struct {
	c *rs.Client
}

func (c *Client) Send(ctx context.Context, m *pb.Mail) error {
	resCh := make(chan *rs.SendEmailResponse)
	errCh := make(chan error)

	go func() {
		res, err := c.c.Emails.Send(&rs.SendEmailRequest{
			From:    m.From,
			To:      m.To,
			Subject: m.Subject,
			Html:    m.Html,
		})
		if err != nil {
			errCh <- err
			return
		}
		resCh <- res
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("%s: failed to send message: %w", operation, err)
	case <-resCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

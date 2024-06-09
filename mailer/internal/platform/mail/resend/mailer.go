package resend

import (
	"context"
	"fmt"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	rs "github.com/resend/resend-go/v2"
)

const operation = "resend mail client"

// NewClient constructs new Resend client
// with provided token.
func NewClient(token string) *Client {
	return &Client{
		c: rs.NewClient(token),
	}
}

// Client is a thin wrapper around resend's SDK
// which will add context support to the existing
// signature call.
type Client struct {
	c *rs.Client
}

// Send method initiates a call to the resend API using
// bult-in resend's SDK. Blocks untill call is finished, or
// error is raised, or context is done.
func (c *Client) Send(ctx context.Context, m *pb.Mail) error {
	if len(m.To) == 0 {
		return fmt.Errorf("%s: recipients cannot be empty", operation)
	}

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

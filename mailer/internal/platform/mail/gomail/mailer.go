package gomail

import (
	"context"
	"fmt"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/mailer"
	"gopkg.in/gomail.v2"
)

const operation = "smtp client"

func NewClient(from, password, host string, port int) *Client {
	d := gomail.NewDialer(host, port, from, password)
	return &Client{
		dialer: d,
		from:   from,
	}
}

type Client struct {
	dialer *gomail.Dialer
	from   string
}

func (c *Client) Send(ctx context.Context, in *pb.Mail) error {
	done := c.send(in)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *Client) send(in *pb.Mail) <-chan error {
	done := make(chan error, 1)

	m := gomail.NewMessage()
	m.SetHeader("From", c.from)
	m.SetHeader("To", in.GetTo()...)
	m.SetHeader("Subject", in.GetSubject())
	m.SetBody("text/html", in.GetHtml())

	go func() {
		if err := c.dialer.DialAndSend(m); err != nil {
			done <- fmt.Errorf("%s: failed to dial and send: %w", operation, err)
			return
		}
		done <- nil
	}()

	return done
}

package mailer

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	pb "github.com/hrvadl/converter/protos/gen/go/v1/mailer"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/converter/sub/internal/transport/grpc/clients/mailer/mocks"
)

func TestClientSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		log  *slog.Logger
		api  pb.MailerServiceClient
		from string
	}
	type args struct {
		ctx     context.Context
		html    string
		subject string
		to      []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, mailer pb.MailerServiceClient)
		wantErr bool
	}{
		{
			name: "Should not return error when was send correctly",
			fields: fields{
				log:  slog.Default(),
				api:  mocks.NewMockMailerServiceClient(gomock.NewController(t)),
				from: "vadym@hrashchenko.com",
			},
			args: args{
				ctx:     context.Background(),
				html:    "test html",
				subject: "test subject",
				to:      []string{"to@to.com", "to1@to.com"},
			},
			setup: func(t *testing.T, mailer pb.MailerServiceClient) {
				m, ok := mailer.(*mocks.MockMailerServiceClient)
				if !ok {
					t.Fatal("Failed to cast mailer client to mock mailer client")
				}

				m.EXPECT().Send(gomock.Any(), &pb.Mail{
					From:    "vadym@hrashchenko.com",
					To:      []string{"to@to.com", "to1@to.com"},
					Subject: "test subject",
					Html:    "test html",
				}).Times(1).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "Should return error when was send incorrectly",
			fields: fields{
				log:  slog.Default(),
				api:  mocks.NewMockMailerServiceClient(gomock.NewController(t)),
				from: "vadym@hrashchenko.com",
			},
			args: args{
				ctx:     context.Background(),
				html:    "test html",
				subject: "test subject",
				to:      []string{"to@to.com", "to1@to.com"},
			},
			setup: func(t *testing.T, mailer pb.MailerServiceClient) {
				m, ok := mailer.(*mocks.MockMailerServiceClient)
				if !ok {
					t.Fatal("Failed to cast mailer client to mock mailer client")
				}

				m.EXPECT().Send(gomock.Any(), &pb.Mail{
					From:    "vadym@hrashchenko.com",
					To:      []string{"to@to.com", "to1@to.com"},
					Subject: "test subject",
					Html:    "test html",
				}).Times(1).Return(nil, errors.New("failed to send"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.api)
			c := &Client{
				log:  tt.fields.log,
				api:  tt.fields.api,
				from: tt.fields.from,
			}
			if err := c.Send(tt.args.ctx, tt.args.html, tt.args.subject, tt.args.to...); (err != nil) != tt.wantErr {
				t.Errorf("Client.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package mailer

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/mailer"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const operation = "mailer server"

func Register(srv *grpc.Server, client Client, log *slog.Logger) {
	pb.RegisterMailerServiceServer(srv, &Server{
		log:    log,
		client: client,
	})
}

//go:generate mockgen -destination=./mocks/mock_client.go -package=mocks . Client
type Client interface {
	Send(ctx context.Context, m *pb.Mail) error
}

type Server struct {
	pb.UnimplementedMailerServiceServer
	log    *slog.Logger
	client Client
}

func (s *Server) Send(ctx context.Context, m *pb.Mail) (*emptypb.Empty, error) {
	if err := s.client.Send(ctx, m); err != nil {
		return nil, fmt.Errorf("%s: failed to send mail: %w", operation, err)
	}
	return nil, nil
}

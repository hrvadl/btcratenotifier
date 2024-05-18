package sub

import (
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/sub"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const operation = "sub server"

func Register(srv *grpc.Server, svc Service, log *slog.Logger) {
	pb.RegisterSubServiceServer(srv, &Server{
		log: log,
		svc: svc,
	})
}

//go:generate mockgen -destination=./mocks/mock_svcr.go -package=mocks . Service
type Service interface {
	Subscribe(ctx context.Context, mail string) (int64, error)
}

type Server struct {
	pb.UnimplementedSubServiceServer
	log *slog.Logger
	svc Service
}

func (s *Server) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (*emptypb.Empty, error) {
	if _, err := s.svc.Subscribe(ctx, req.Email); err != nil {
		return nil, fmt.Errorf("%s: failed to subscribe user: %w", operation, err)
	}
	return nil, nil
}

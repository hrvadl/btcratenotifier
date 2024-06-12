package sub

import (
	"errors"
	"log/slog"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	subSvc "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sub"
)

// Register registers subscribe handler to the given GRPC server.
// NOTE: all parameters are required, the service will panic if
// either of them is missing.
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

// Server represents subscribe GRPC server
// which will handle the incoming requests and delegate
// all work to the underlying svc.
type Server struct {
	pb.UnimplementedSubServiceServer
	log *slog.Logger
	svc Service
}

// Subscribe method calls underlying service method and returns an error, in case there was a
// failure.
func (s *Server) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (*emptypb.Empty, error) {
	_, err := s.svc.Subscribe(ctx, req.GetEmail())
	if err == nil {
		return &emptypb.Empty{}, nil
	}

	if errors.Is(err, subSvc.ErrInvalidEmail) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, subSvc.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	return nil, status.Error(codes.Internal, err.Error())
}

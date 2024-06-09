package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const operation = "converted server"

// Registers rate watcher handler to the given GRPC server.
// NOTE: all parameters are required, the service will panic if
// either of them is missing.
func Register(srv *grpc.Server, cnv Converter, log *slog.Logger) {
	pb.RegisterRateWatcherServiceServer(srv, &Server{
		log:       log,
		converter: cnv,
	})
}

//go:generate mockgen -destination=./mocks/mock_converter.go -package=mocks . Converter
type Converter interface {
	Convert(ctx context.Context) (float32, error)
}

// Server represents rate watcher GRPC server
// which will handle the incoming requests and delegate
// all work to the underlying converter.
type Server struct {
	pb.UnimplementedRateWatcherServiceServer
	log       *slog.Logger
	converter Converter
}

// GetRate method calls underlying converter method and returns an error, in case there was a
// failure.
func (s *Server) GetRate(ctx context.Context, _ *emptypb.Empty) (*pb.RateResponse, error) {
	rate, err := s.converter.Convert(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to convert: %w", operation, err)
	}
	return &pb.RateResponse{Rate: rate}, nil
}

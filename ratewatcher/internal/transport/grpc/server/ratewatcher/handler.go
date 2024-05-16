package ratewatcher

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const operation = "converted server"

func Register(srv *grpc.Server, cnv Converter, log *slog.Logger) {
	pb.RegisterRateWatcherServiceServer(srv, &Server{
		log:       log,
		converter: cnv,
	})
}

type Converter interface {
	Convert(ctx context.Context) (float32, error)
}

type Server struct {
	pb.UnimplementedRateWatcherServiceServer
	log       *slog.Logger
	converter Converter
}

func (s *Server) GetRate(ctx context.Context, _ *emptypb.Empty) (*pb.RateResponse, error) {
	rate, err := s.converter.Convert(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to convert: %w", operation, err)
	}
	return &pb.RateResponse{Rate: rate}, nil
}

//go:build !integration

package ratewatcher

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/rw/internal/transport/grpc/server/ratewatcher/mocks"
)

func TestServerGetRate(t *testing.T) {
	t.Parallel()
	type fields struct {
		log       *slog.Logger
		converter Converter
	}
	type args struct {
		ctx context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, converter Converter)
		want    *pb.RateResponse
		wantErr bool
	}{
		{
			name: "Should not return error when converter succeeded",
			fields: fields{
				log:       slog.Default(),
				converter: mocks.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				in1: nil,
			},
			setup: func(t *testing.T, converter Converter) {
				t.Helper()
				c, ok := converter.(*mocks.MockConverter)
				require.True(t, ok, "Failed to cast converter to mock converter")
				c.EXPECT().Convert(gomock.Any()).Times(1).Return(float32(3.3), nil)
			},
			want:    &pb.RateResponse{Rate: 3.3},
			wantErr: false,
		},
		{
			name: "Should return error when converter failed",
			fields: fields{
				log:       slog.Default(),
				converter: mocks.NewMockConverter(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				in1: nil,
			},
			setup: func(t *testing.T, converter Converter) {
				t.Helper()
				c, ok := converter.(*mocks.MockConverter)
				require.True(t, ok, "Failed to cast converter to mock converter")
				c.EXPECT().
					Convert(gomock.Any()).
					Times(1).
					Return(float32(0), errors.New("failed to convert USD to UAH"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.converter)
			s := &Server{
				UnimplementedRateWatcherServiceServer: pb.UnimplementedRateWatcherServiceServer{},
				log:                                   tt.fields.log,
				converter:                             tt.fields.converter,
			}

			got, err := s.GetRate(tt.args.ctx, tt.args.in1)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

//go:generate mockgen -destination=./mocks/mock_rw.go -package=mocks github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher RateWatcherServiceClient
package ratewatcher

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/transport/grpc/clients/ratewatcher/mocks"
)

func TestClientGetRate(t *testing.T) {
	t.Parallel()
	type fields struct {
		log *slog.Logger
		api pb.RateWatcherServiceClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, rws pb.RateWatcherServiceClient)
		want    float32
		wantErr bool
	}{
		{
			name: "Should not return error when rate service succeeded",
			fields: fields{
				log: slog.Default(),
				api: mocks.NewMockRateWatcherServiceClient(gomock.NewController(t)),
			},
			args: args{ctx: context.Background()},
			setup: func(t *testing.T, rws pb.RateWatcherServiceClient) {
				t.Helper()
				rw, ok := rws.(*mocks.MockRateWatcherServiceClient)
				require.True(t, ok, "Failed to cast rw to mock rw")
				rw.EXPECT().
					GetRate(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&pb.RateResponse{Rate: 1}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Should return error when rate service failed",
			fields: fields{
				log: slog.Default(),
				api: mocks.NewMockRateWatcherServiceClient(gomock.NewController(t)),
			},
			args: args{ctx: context.Background()},
			setup: func(t *testing.T, rws pb.RateWatcherServiceClient) {
				t.Helper()
				rw, ok := rws.(*mocks.MockRateWatcherServiceClient)
				require.True(t, ok, "Failed to cast rw to mock rw")
				rw.EXPECT().
					GetRate(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("failed to get rate"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.api)
			c := &Client{
				log: tt.fields.log,
				api: tt.fields.api,
			}

			got, err := c.GetRate(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.InEpsilon(t, tt.want, got, 2)
		})
	}
}

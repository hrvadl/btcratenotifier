//go:generate mockgen -destination=./mocks/mock_rw.go -package=mocks github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher RateWatcherServiceClient
package ratewatcher

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/btcratenotifier/sub/internal/transport/grpc/clients/ratewatcher/mocks"
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
				rw, ok := rws.(*mocks.MockRateWatcherServiceClient)
				if !ok {
					t.Fatal("Failed to cast rw to mock rw")
				}

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
				rw, ok := rws.(*mocks.MockRateWatcherServiceClient)
				if !ok {
					t.Fatal("Failed to cast rw to mock rw")
				}

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
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Client.GetRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

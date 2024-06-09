package ratewatcher

import (
	"context"
	"errors"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/ratewatcher"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/grpc/clients/ratewatcher/mocks"
)

func TestClientGetRate(t *testing.T) {
	t.Parallel()
	type fields struct {
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
			name: "Should not return error when rate watcher svc succeeded",
			fields: fields{
				api: mocks.NewMockRateWatcherServiceClient(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, rws pb.RateWatcherServiceClient) {
				t.Helper()
				rw, ok := rws.(*mocks.MockRateWatcherServiceClient)
				if !ok {
					t.Fatal("failed to cast rate watcher client to mock")
				}

				rw.EXPECT().
					GetRate(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&pb.RateResponse{Rate: float32(39.3)}, nil)
			},
			want:    39.3,
			wantErr: false,
		},
		{
			name: "Should not return error when rate watcher svc succeeded",
			fields: fields{
				api: mocks.NewMockRateWatcherServiceClient(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
			},
			setup: func(t *testing.T, rws pb.RateWatcherServiceClient) {
				t.Helper()
				rw, ok := rws.(*mocks.MockRateWatcherServiceClient)
				if !ok {
					t.Fatal("failed to cast rate watcher client to mock")
				}

				rw.EXPECT().
					GetRate(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("failed to get exchange rate"))
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

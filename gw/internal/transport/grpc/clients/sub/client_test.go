//go:generate mockgen -destination=./mocks/mock_sub.go -package=mocks github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub SubServiceClient
package sub

import (
	"context"
	"errors"
	"testing"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/grpc/clients/sub/mocks"
)

func TestClientSubscribe(t *testing.T) {
	t.Parallel()
	type fields struct {
		api pb.SubServiceClient
	}
	type args struct {
		ctx context.Context
		req *pb.SubscribeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, subscriber pb.SubServiceClient)
		wantErr bool
	}{
		{
			name: "Should not return error when subscribe svc succeeded",
			fields: fields{
				api: mocks.NewMockSubServiceClient(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				req: &pb.SubscribeRequest{
					Email: "sub@me.com",
				},
			},
			setup: func(t *testing.T, subscriber pb.SubServiceClient) {
				t.Helper()
				s, ok := subscriber.(*mocks.MockSubServiceClient)
				require.True(t, ok, "Failed to convert subscribe server client")
				s.EXPECT().Subscribe(gomock.Any(), &pb.SubscribeRequest{
					Email: "sub@me.com",
				}).Times(1).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "Should return error when subscribe svc failed",
			fields: fields{
				api: mocks.NewMockSubServiceClient(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				req: &pb.SubscribeRequest{
					Email: "sub@me.com",
				},
			},
			setup: func(t *testing.T, subscriber pb.SubServiceClient) {
				t.Helper()
				s, ok := subscriber.(*mocks.MockSubServiceClient)
				require.True(t, ok, "Failed to convert subscribe server client")
				s.EXPECT().Subscribe(gomock.Any(), &pb.SubscribeRequest{
					Email: "sub@me.com",
				}).Times(1).Return(nil, errors.New("failed to subscribe user"))
			},
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

			err := c.Subscribe(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

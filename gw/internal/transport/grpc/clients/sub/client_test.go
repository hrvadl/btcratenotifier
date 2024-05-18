//go:generate mockgen -destination=./mocks/mock_sub.go -package=mocks github.com/hrvadl/converter/protos/gen/go/v1/sub SubServiceClient
package sub

import (
	"context"
	"errors"
	"testing"

	pb "github.com/hrvadl/converter/protos/gen/go/v1/sub"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/converter/gw/internal/transport/grpc/clients/sub/mocks"
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
				s, ok := subscriber.(*mocks.MockSubServiceClient)
				if !ok {
					t.Fatal("Failed to convert subscribe server client")
				}

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
				s, ok := subscriber.(*mocks.MockSubServiceClient)
				if !ok {
					t.Fatal("Failed to convert subscribe server client")
				}

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
			if err := c.Subscribe(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Client.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

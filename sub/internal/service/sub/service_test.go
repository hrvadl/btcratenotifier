package sub

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/hrvadl/converter/sub/internal/service/sub/mocks"
	"github.com/hrvadl/converter/sub/internal/storage/subscriber"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	type args struct {
		rr RecipientSaver
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Should create new service correctly when correct arguments are provided",
			args: args{
				rr: mocks.NewMockRecipientSaver(gomock.NewController(t)),
			},
			want: &Service{
				repo: mocks.NewMockRecipientSaver(gomock.NewController(t)),
			},
		},
		{
			name: "Should create new service correctly when allowed arguments are provided",
			args: args{
				rr: nil,
			},
			want: &Service{
				repo: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewService(tt.args.rr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceSubscribe(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo RecipientSaver
	}
	type args struct {
		ctx  context.Context
		mail string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
		setup   func(t *testing.T, saver RecipientSaver)
	}{
		{
			name: "Should not return err when everything is correct",
			fields: fields{
				repo: mocks.NewMockRecipientSaver(gomock.NewController(t)),
			},
			args: args{
				ctx:  context.Background(),
				mail: "mail@gmail.com",
			},
			setup: func(t *testing.T, saver RecipientSaver) {
				rs, ok := saver.(*mocks.MockRecipientSaver)
				if !ok {
					t.Fatalf("Failed to cast saver to mock saver")
				}
				rs.EXPECT().
					Save(gomock.Any(), subscriber.Subscriber{Email: "mail@gmail.com"}).
					Times(1).
					Return(int64(1), nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Should return err when saver returned err",
			fields: fields{
				repo: mocks.NewMockRecipientSaver(gomock.NewController(t)),
			},
			args: args{
				ctx:  context.Background(),
				mail: "mail@gmail.com",
			},
			setup: func(t *testing.T, saver RecipientSaver) {
				rs, ok := saver.(*mocks.MockRecipientSaver)
				if !ok {
					t.Fatalf("Failed to cast saver to mock saver")
				}
				rs.EXPECT().
					Save(gomock.Any(), subscriber.Subscriber{Email: "mail@gmail.com"}).
					Times(1).
					Return(int64(0), errors.New("failed to save subscriber"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Should return err when email is empty",
			fields: fields{
				repo: mocks.NewMockRecipientSaver(gomock.NewController(t)),
			},
			args: args{
				ctx:  context.Background(),
				mail: "",
			},
			setup: func(t *testing.T, saver RecipientSaver) {
				rs, ok := saver.(*mocks.MockRecipientSaver)
				if !ok {
					t.Fatalf("Failed to cast saver to mock saver")
				}
				rs.EXPECT().
					Save(gomock.Any(), subscriber.Subscriber{Email: "mail@gmail.com"}).
					Times(0).
					Return(int64(0), errors.New("failed to save subscriber"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.repo)
			s := &Service{repo: tt.fields.repo}
			got, err := s.Subscribe(tt.args.ctx, tt.args.mail)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Service.Subscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}

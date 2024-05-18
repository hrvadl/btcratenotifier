package sender

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/hrvadl/btcratenotifier/sub/internal/service/sender/mocks"
)

func TestCronJobAdapterDo(t *testing.T) {
	t.Parallel()
	type fields struct {
		sender Sender
		log    *slog.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(t *testing.T, m Sender)
		wantErr bool
	}{
		{
			name: "Should return err when sender failed",
			fields: fields{
				sender: mocks.NewMockSender(gomock.NewController(t)),
				log:    slog.Default(),
			},
			wantErr: true,
			setup: func(t *testing.T, m Sender) {
				t.Helper()
				ss, ok := m.(*mocks.MockSender)
				if !ok {
					t.Fatal("failed to cast sender to mock")
				}
				err := errors.New("failed to send")
				ss.EXPECT().Send(gomock.Any()).Times(1).Return(err)
			},
		},
		{
			name: "Should not return err when sender succeeded",
			fields: fields{
				sender: mocks.NewMockSender(gomock.NewController(t)),
				log:    slog.Default(),
			},
			wantErr: false,
			setup: func(t *testing.T, m Sender) {
				t.Helper()
				ss, ok := m.(*mocks.MockSender)
				if !ok {
					t.Fatal("failed to cast sender to mock")
				}
				ss.EXPECT().Send(gomock.Any()).Times(1).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.sender)
			c := &CronJobAdapter{
				sender: tt.fields.sender,
				log:    tt.fields.log,
			}

			if err := c.Do(); (err != nil) != tt.wantErr {
				t.Errorf("CronJobAdapter.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCronJobAdapter(t *testing.T) {
	t.Parallel()
	type args struct {
		s   Sender
		log *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *CronJobAdapter
	}{
		{
			name: "Should return crob job adapter with correct arguments provided",
			args: args{
				s:   mocks.NewMockSender(gomock.NewController(t)),
				log: slog.Default(),
			},
			want: &CronJobAdapter{
				sender: mocks.NewMockSender(gomock.NewController(t)),
				log:    slog.Default(),
			},
		},
		{
			name: "Should return crob job adapter with allowed arguments provided",
			args: args{
				s:   nil,
				log: nil,
			},
			want: &CronJobAdapter{
				sender: nil,
				log:    nil,
			},
		},
		{
			name: "Should return crob job adapter with allowed arguments provided",
			args: args{
				s:   nil,
				log: slog.Default(),
			},
			want: &CronJobAdapter{
				sender: nil,
				log:    slog.Default(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewCronJobAdapter(tt.args.s, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCronJobAdapter() = %v, want %v", got, tt.want)
			}
		})
	}
}

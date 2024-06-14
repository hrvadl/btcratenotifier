package sender

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sender/mocks"
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
				require.True(t, ok, "Failed to cast sender to mock")
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
				require.True(t, ok, "Failed to cast sender to mock")
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

			err := c.Do()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewCronJobAdapter(t *testing.T) {
	t.Parallel()
	type args struct {
		s       Sender
		log     *slog.Logger
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *CronJobAdapter
	}{
		{
			name: "Should return crob job adapter with correct arguments provided",
			args: args{
				s:       mocks.NewMockSender(gomock.NewController(t)),
				log:     slog.Default(),
				timeout: time.Second,
			},
			want: &CronJobAdapter{
				sender:  mocks.NewMockSender(gomock.NewController(t)),
				log:     slog.Default(),
				timeout: time.Second,
			},
		},
		{
			name: "Should return crob job adapter with correct arguments provided",
			args: args{
				s:       mocks.NewMockSender(gomock.NewController(t)),
				log:     slog.Default(),
				timeout: time.Microsecond,
			},
			want: &CronJobAdapter{
				sender:  mocks.NewMockSender(gomock.NewController(t)),
				log:     slog.Default(),
				timeout: time.Microsecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCronJobAdapter(tt.args.s, tt.args.timeout, tt.args.log)
			require.Equal(t, tt.want, got)
		})
	}
}

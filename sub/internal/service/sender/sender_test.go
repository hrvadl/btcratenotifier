package sender

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sender/mocks"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/subscriber"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		m   Mailer
		sg  SubscriberGetter
		mf  RateMessageFormatter
		rg  RateGetter
		log *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Shoild initialize sender service when correct arguments are provided",
			args: args{
				m:   mocks.NewMockMailer(gomock.NewController(t)),
				sg:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				mf:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rg:  mocks.NewMockRateGetter(gomock.NewController(t)),
				log: slog.Default(),
			},
			want: &Service{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: mocks.NewMockRateGetter(gomock.NewController(t)),
				log:        slog.Default(),
			},
		},
		{
			name: "Shoild initialize sender service when allowed arguments are provided",
			args: args{
				m:   mocks.NewMockMailer(gomock.NewController(t)),
				sg:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				mf:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rg:  nil,
				log: nil,
			},
			want: &Service{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: nil,
				log:        nil,
			},
		},
		{
			name: "Shoild initialize sender service when nil arguments are provided",
			args: args{
				m:   nil,
				sg:  nil,
				mf:  nil,
				rg:  nil,
				log: nil,
			},
			want: &Service{
				mailer:     nil,
				subGetter:  nil,
				formatter:  nil,
				rateGetter: nil,
				log:        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New(tt.args.m, tt.args.sg, tt.args.mf, tt.args.rg, tt.args.log)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServiceSend(t *testing.T) {
	t.Parallel()
	type fields struct {
		mailer     Mailer
		formatter  RateMessageFormatter
		subGetter  SubscriberGetter
		rateGetter RateGetter
		log        *slog.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(t *testing.T, f *fields)
	}{
		{
			name: "Should not return error when everything is correct",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: mocks.NewMockRateGetter(gomock.NewController(t)),
				log:        slog.Default(),
			},
			setup: func(t *testing.T, f *fields) {
				t.Helper()
				var (
					rate   float32 = 10.
					fmtMsg         = "fmtTestMsg"
					subs           = []subscriber.Subscriber{
						{ID: 1, Email: "test@test.com"},
						{ID: 2, Email: "test2@test.com"},
					}
				)

				m, ok := f.mailer.(*mocks.MockMailer)
				require.True(t, ok, "failed to cast mailer to mock mailer")
				m.EXPECT().
					Send(gomock.Any(), fmtMsg, subject, "test@test.com", "test2@test.com").
					Times(1).
					Return(nil)

				rg, ok := f.rateGetter.(*mocks.MockRateGetter)
				require.True(t, ok, "failed to cast getter to mock getter")
				rg.EXPECT().GetRate(gomock.Any()).Times(1).Return(rate, nil)

				fmter, ok := f.formatter.(*mocks.MockRateMessageFormatter)
				require.True(t, ok, "failed to cast fmter to mock fmter")
				fmter.EXPECT().Format(rate).Times(1).Return(fmtMsg)

				sg, ok := f.subGetter.(*mocks.MockSubscriberGetter)
				require.True(t, ok, "failed to cast sub getter to mock sub getter")
				sg.EXPECT().Get(gomock.Any()).Times(1).Return(subs, nil)
			},
			wantErr: false,
		},
		{
			name: "Should return error when subs getter returned err",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: mocks.NewMockRateGetter(gomock.NewController(t)),
				log:        slog.Default(),
			},
			setup: func(t *testing.T, f *fields) {
				t.Helper()
				var (
					rate   float32 = 10.
					fmtMsg         = "fmtTestMsg"
					subs           = []subscriber.Subscriber{
						{ID: 1, Email: "test@test.com"},
						{ID: 2, Email: "test2@test.com"},
					}
				)

				m, ok := f.mailer.(*mocks.MockMailer)
				require.True(t, ok, "failed to cast mailer to mock mailer")
				m.EXPECT().
					Send(gomock.Any(), fmtMsg, subject, "test@test.com", "test2@test.com").
					Times(0).
					Return(nil)

				rg, ok := f.rateGetter.(*mocks.MockRateGetter)
				require.True(t, ok, "failed to cast getter to mock getter")
				rg.EXPECT().GetRate(gomock.Any()).Times(0).Return(rate, nil)

				fmter, ok := f.formatter.(*mocks.MockRateMessageFormatter)
				require.True(t, ok, "failed to cast fmter to mock fmter")
				fmter.EXPECT().Format(rate).Times(0).Return(fmtMsg)

				sg, ok := f.subGetter.(*mocks.MockSubscriberGetter)
				require.True(t, ok, "failed to cast sub getter to mock sub getter")
				fmter.EXPECT().Format(rate).Times(0).Return(fmtMsg)
				sg.EXPECT().
					Get(gomock.Any()).
					Times(1).
					Return(subs, errors.New("failed to get subs"))
			},
			wantErr: true,
		},
		{
			name: "Should return error when subs are empty",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: mocks.NewMockRateGetter(gomock.NewController(t)),
				log:        slog.Default(),
			},
			setup: func(t *testing.T, f *fields) {
				t.Helper()
				var (
					rate   float32 = 10.
					fmtMsg         = "fmtTestMsg"
				)

				m, ok := f.mailer.(*mocks.MockMailer)
				if !ok {
					t.Fatal("failed to cast mailer to mock mailer")
				}

				m.EXPECT().
					Send(gomock.Any(), fmtMsg, subject, "test@test.com", "test2@test.com").
					Times(0).
					Return(nil)

				rg, ok := f.rateGetter.(*mocks.MockRateGetter)
				require.True(t, ok, "failed to cast getter to mock getter")
				rg.EXPECT().GetRate(gomock.Any()).Times(0).Return(rate, nil)

				fmter, ok := f.formatter.(*mocks.MockRateMessageFormatter)
				require.True(t, ok, "failed to cast fmter to mock fmter")
				fmter.EXPECT().Format(rate).Times(0).Return(fmtMsg)

				sg, ok := f.subGetter.(*mocks.MockSubscriberGetter)
				require.True(t, ok, "failed to cast sub getter to mock sub getter")
				sg.EXPECT().Get(gomock.Any()).Times(1).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "Should return error when rate getter returned err",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: mocks.NewMockRateGetter(gomock.NewController(t)),
				log:        slog.Default(),
			},
			setup: func(t *testing.T, f *fields) {
				t.Helper()
				var (
					rate   float32 = 10.
					fmtMsg         = "fmtTestMsg"
					subs           = []subscriber.Subscriber{
						{ID: 1, Email: "test@test.com"},
						{ID: 2, Email: "test2@test.com"},
					}
				)

				m, ok := f.mailer.(*mocks.MockMailer)
				require.True(t, ok, "failed to cast mailer to mock mailer")
				m.EXPECT().
					Send(gomock.Any(), fmtMsg, subject, "test@test.com", "test2@test.com").
					Times(0).
					Return(nil)

				rg, ok := f.rateGetter.(*mocks.MockRateGetter)
				require.True(t, ok, "failed to cast rate getter to mock getter")
				rg.EXPECT().
					GetRate(gomock.Any()).
					Times(1).
					Return(rate, errors.New("failed to get rate"))

				fmter, ok := f.formatter.(*mocks.MockRateMessageFormatter)
				require.True(t, ok, "failed to cast fmter to mock fmter")
				fmter.EXPECT().Format(rate).Times(0).Return(fmtMsg)

				sg, ok := f.subGetter.(*mocks.MockSubscriberGetter)
				require.True(t, ok, "failed to cast sub getter to mock sub getter")
				sg.EXPECT().
					Get(gomock.Any()).
					Times(1).
					Return(subs, nil)
			},
			wantErr: true,
		},
		{
			name: "Should return error when mailer returned err",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				mailer:     mocks.NewMockMailer(gomock.NewController(t)),
				subGetter:  mocks.NewMockSubscriberGetter(gomock.NewController(t)),
				formatter:  mocks.NewMockRateMessageFormatter(gomock.NewController(t)),
				rateGetter: mocks.NewMockRateGetter(gomock.NewController(t)),
				log:        slog.Default(),
			},
			setup: func(t *testing.T, f *fields) {
				t.Helper()
				var (
					rate   float32 = 10.
					fmtMsg         = "fmtTestMsg"
					subs           = []subscriber.Subscriber{
						{ID: 1, Email: "test@test.com"},
						{ID: 2, Email: "test2@test.com"},
					}
				)

				m, ok := f.mailer.(*mocks.MockMailer)
				require.True(t, ok, "failed to cast mailer to mock mailer")
				m.EXPECT().
					Send(gomock.Any(), fmtMsg, subject, "test@test.com", "test2@test.com").
					Times(1).
					Return(errors.New("failed to send msg"))

				rg, ok := f.rateGetter.(*mocks.MockRateGetter)
				require.True(t, ok, "failed to cast rate getter to mock rate getter")
				rg.EXPECT().
					GetRate(gomock.Any()).
					Times(1).
					Return(rate, nil)

				fmter, ok := f.formatter.(*mocks.MockRateMessageFormatter)
				require.True(t, ok, "failed to cast fmter to mock fmter")
				fmter.EXPECT().Format(rate).Times(1).Return(fmtMsg)

				sg, ok := f.subGetter.(*mocks.MockSubscriberGetter)
				require.True(t, ok, "failed to cast sub getter to mock sub getter")
				sg.EXPECT().
					Get(gomock.Any()).
					Times(1).
					Return(subs, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, &tt.fields)
			w := &Service{
				mailer:     tt.fields.mailer,
				formatter:  tt.fields.formatter,
				subGetter:  tt.fields.subGetter,
				rateGetter: tt.fields.rateGetter,
				log:        tt.fields.log,
			}

			err := w.Send(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestMapSubsToMails(t *testing.T) {
	t.Parallel()
	type args struct {
		s []subscriber.Subscriber
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Should map subscribers to email correctly",
			args: args{
				s: []subscriber.Subscriber{
					{Email: "test@test.com"},
					{Email: "test2@test.com"},
					{Email: "test3@test.com"},
					{Email: "test4@test.com"},
				},
			},
			want: []string{
				"test@test.com",
				"test2@test.com",
				"test3@test.com",
				"test4@test.com",
			},
		},
		{
			name: "Should map empty subscribers to email correctly",
			args: args{
				s: []subscriber.Subscriber{},
			},
			want: []string{},
		},
		{
			name: "Should map nil subscribers to email correctly",
			args: args{
				s: nil,
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mapSubsToMails(tt.args.s)
			require.Equal(t, tt.want, got)
		})
	}
}

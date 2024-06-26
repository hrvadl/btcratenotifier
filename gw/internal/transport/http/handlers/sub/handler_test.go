package sub

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	pb "github.com/hrvadl/converter/protos/gen/go/v1/sub"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/converter/gw/internal/transport/http/handlers/sub/mocks"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		svc Service
		log *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "Should create handler correctly",
			args: args{
				log: slog.Default(),
				svc: mocks.NewMockService(gomock.NewController(t)),
			},
			want: &Handler{
				log: slog.Default(),
				svc: mocks.NewMockService(gomock.NewController(t)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewHandler(tt.args.svc, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerSubscribe(t *testing.T) {
	t.Parallel()
	type fields struct {
		svc Service
		log *slog.Logger
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, service Service)
		want   int
	}{
		{
			name: "Should return 200 when service succeeded",
			fields: fields{
				svc: mocks.NewMockService(gomock.NewController(t)),
				log: slog.Default(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: withFormDataContentType(httptest.NewRequest(
					http.MethodPost,
					"/",
					bytes.NewBufferString(url.Values{"email": {"test@test.com"}}.Encode()),
				)),
			},
			setup: func(t *testing.T, service Service) {
				t.Helper()
				svc, ok := service.(*mocks.MockService)
				if !ok {
					t.Fatal("Failed to cast service to mock")
				}

				svc.EXPECT().
					Subscribe(gomock.Any(), &pb.SubscribeRequest{Email: "test@test.com"}).
					Times(1).
					Return(nil)
			},
			want: http.StatusOK,
		},
		{
			name: "Should return 409 when service failed",
			fields: fields{
				svc: mocks.NewMockService(gomock.NewController(t)),
				log: slog.Default(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: withFormDataContentType(httptest.NewRequest(
					http.MethodPost,
					"/",
					bytes.NewBufferString(url.Values{"email": {"test@test.com"}}.Encode()),
				)),
			},
			setup: func(t *testing.T, service Service) {
				t.Helper()
				svc, ok := service.(*mocks.MockService)
				if !ok {
					t.Fatal("Failed to cast service to mock")
				}

				svc.EXPECT().
					Subscribe(gomock.Any(), &pb.SubscribeRequest{Email: "test@test.com"}).
					Times(1).
					Return(errors.New("failed to subscribe"))
			},
			want: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, tt.fields.svc)
			h := &Handler{
				svc: tt.fields.svc,
				log: tt.fields.log,
			}
			h.Subscribe(tt.args.w, tt.args.r)
			if got := tt.args.w.Result().StatusCode; got != tt.want {
				t.Errorf("Subscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func withFormDataContentType(r *http.Request) *http.Request {
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

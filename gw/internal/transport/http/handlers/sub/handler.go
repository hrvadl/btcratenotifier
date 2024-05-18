package sub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	pb "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/sub"
)

func NewHandler(svc Service, log *slog.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

//go:generate mockgen -destination=./mocks/mock_svc.go -package=mocks . Service
type Service interface {
	Subscribe(ctx context.Context, s *pb.SubscribeRequest) error
}

type Handler struct {
	svc Service
	log *slog.Logger
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	f := r.Form
	h.log.Info(fmt.Sprintf("%v", f))
	mail := r.FormValue("email")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if err := h.svc.Subscribe(ctx, &pb.SubscribeRequest{Email: mail}); err != nil {
		h.log.Error("Failed to subscribe user", "err", err)
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte("Email is already subscribed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Added email."))
}

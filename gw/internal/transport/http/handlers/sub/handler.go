package sub

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/http/handlers"
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

// Subscribe godoc
// @Summary      Subscribe to email rate exchange notification
// @Tags         Rate
// @Accept       application/x-www-form-urlencoded
// @Produce      json
// @Param        body formData string true "Email to subscribe"
// @Success      200  {object}  handlers.EmptyResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Router       /api/subscribe [post]
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	mail := r.FormValue("email")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if err := h.svc.Subscribe(ctx, &pb.SubscribeRequest{Email: mail}); err != nil {
		h.log.Error("Failed to subscribe user", "err", err)
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write(handlers.NewErrResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(handlers.NewEmptyResponse("added email"))
}

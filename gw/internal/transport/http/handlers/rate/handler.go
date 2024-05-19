package rate

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hrvadl/converter/gw/internal/transport/http/handlers"
)

func NewHandler(rg Getter, log *slog.Logger) *Handler {
	return &Handler{
		rg:  rg,
		log: log,
	}
}

//go:generate mockgen -destination=./mocks/mock_getter.go -package=mocks . Getter
type Getter interface {
	GetRate(ctx context.Context) (float32, error)
}

type Handler struct {
	log *slog.Logger
	rg  Getter
}

func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	rat, err := h.rg.GetRate(ctx)
	if err != nil {
		h.log.Error("Failed to get rate", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(handlers.NewErrResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(handlers.NewSuccessResponse("successfully got rate", rat))
}

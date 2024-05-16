package rate

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func NewHandler(rg Getter) *Handler {
	return &Handler{
		rg: rg,
	}
}

type Getter interface {
	GetRate(ctx context.Context) (float32, error)
}

type Handler struct {
	rg Getter
}

func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	rat, err := h.rg.GetRate(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid status value"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(rat)))
}

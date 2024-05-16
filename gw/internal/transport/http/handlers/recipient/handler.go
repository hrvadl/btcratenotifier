package recipient

import (
	"context"
	"net/http"
	"time"
)

type Service interface {
	Subscribe(ctx context.Context, mail string) error
	SendToAll(ctx context.Context) error
}

func New(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

type Handler struct {
	svc Service
}

func (h *Handler) AddRecipient(w http.ResponseWriter, r *http.Request) {
	mail := r.FormValue("email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.svc.Subscribe(ctx, mail); err != nil {
		w.Write([]byte("E-mail вже є в базі"))
		w.WriteHeader(http.StatusConflict)
	}

	w.Write([]byte("E-mail додано"))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SendToAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := h.svc.SendToAll(ctx); err != nil {
		w.Write([]byte("Поштові листи не були відправлені"))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte("Листи успішно відправлено"))
	w.WriteHeader(http.StatusOK)
}

package sender

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Service interface {
	Subscribe(ctx context.Context, mail string) error
	SendToAll(ctx context.Context) error
}

func NewHandler(svc Service, log *slog.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

type Handler struct {
	svc Service
	log *slog.Logger
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	mail := r.FormValue("email")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if err := h.svc.Subscribe(ctx, mail); err != nil {
		h.log.Error("Failed to subscribe user", "err", err)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Email is already subscribed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Added email."))
}

func (h *Handler) SendToAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*25)
	defer cancel()

	if err := h.svc.SendToAll(ctx); err != nil {
		h.log.Error("Failed to send mailing list", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to send mailing list."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mailing lists were successfuly sent!"))
}

package sub

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	pb "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/protos/gen/go/v1/sub"

	subSvc "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/grpc/clients/sub"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/gw/internal/transport/http/handlers"
)

const subscribeTimeout = 5 * time.Second

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
// @Failure      409  {object}  handlers.ErrorResponse
// @Router       /api/subscribe [post]
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), subscribeTimeout)
	defer cancel()

	mail := r.FormValue("email")
	err := h.svc.Subscribe(ctx, &pb.SubscribeRequest{Email: mail})
	if err == nil {
		h.httpSuccess(w, http.StatusOK, nil, "aded email")
		return
	}

	if errors.Is(err, subSvc.ErrAlreadyExists) {
		h.httpFail(w, http.StatusConflict, err)
	}

	if errors.Is(err, subSvc.ErrInvalidEmail) {
		h.httpFail(w, http.StatusBadRequest, err)
	}

	h.httpFail(w, http.StatusInternalServerError, err)
}

func (h *Handler) httpSuccess(w http.ResponseWriter, code int, data any, msg string) {
	res, err := handlers.NewSuccessResponse(msg, data)
	if err != nil {
		h.log.Error("Failed to construct success response", slog.Any("err", err))
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(res); err != nil {
		h.log.Error("Failed to write response", slog.Any("err", err))
	}
}

func (h *Handler) httpFail(w http.ResponseWriter, code int, err error) {
	res, err := handlers.NewErrResponse(err)
	if err != nil {
		h.log.Error("Failed to construct error response", slog.Any("err", err))
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(res); err != nil {
		h.log.Error("Failed to write response", slog.Any("err", err))
	}
}

package updateStatus

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "projectsShowcase/internal/lib/api/response"
	"projectsShowcase/internal/lib/logger/sl"
	"strconv"
)

type Request struct {
	Status string `json:"status" validate:"required"`
}

type Response struct {
	resp.Response
}

type ApplicationStatusUpdater interface {
	UpdateApplicationStatus(id int64, status string) error
}

func New(log *slog.Logger, applicationStatusUpdater ApplicationStatusUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.application.updateStatus.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid ID format", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid ID format"))
			return
		}

		var req Request

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		err = applicationStatusUpdater.UpdateApplicationStatus(id, req.Status)
		if err != nil {
			log.Error("failed to update application", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to update application"))

			return
		}

		log.Info("application updated", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}

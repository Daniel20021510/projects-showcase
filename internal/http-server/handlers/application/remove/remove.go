package remove

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "projectsShowcase/internal/lib/api/response"
	"projectsShowcase/internal/lib/logger/sl"
	"projectsShowcase/internal/storage"
	"strconv"
)

type ApplicationRemover interface {
	DeleteApplication(id int64) error
}

func New(log *slog.Logger, applicationRemover ApplicationRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.application.remove.New"

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

		err = applicationRemover.DeleteApplication(id)
		if err != nil {
			if errors.Is(err, storage.ErrApplicationNotFound) {
				log.Info("application not found", slog.Int64("id", id))
				render.JSON(w, r, resp.Error("application not found"))
				return
			}
			log.Error("failed to delete application", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete application"))
			return
		}

		log.Info("application deleted", slog.Int64("id", id))
		render.JSON(w, r, resp.OK())
	}
}

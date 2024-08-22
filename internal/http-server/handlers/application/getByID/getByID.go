package getByID

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"projectsShowcase/internal/domain/models"
	resp "projectsShowcase/internal/lib/api/response"
	"projectsShowcase/internal/lib/logger/sl"
)

type Response struct {
	resp.Response
	Application *models.Application `json:"application,omitempty"`
}

type ApplicationGetter interface {
	GetApplicationByID(id string) (*models.Application, error)
}

func New(log *slog.Logger, applicationGetter ApplicationGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.application.getByID.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")

		application, err := applicationGetter.GetApplicationByID(id)
		if err != nil {
			log.Error("failed to get application by ID", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get application by ID"))

			return
		}

		log.Info("get application by ID", slog.String("id", id))

		responseOK(w, r, application)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, application *models.Application) {
	render.JSON(w, r, Response{
		Response:    resp.OK(),
		Application: application,
	})
}

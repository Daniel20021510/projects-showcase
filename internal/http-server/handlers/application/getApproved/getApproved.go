package getApproved

import (
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
	Applications []models.Application `json:"applications,omitempty"`
}

type ApprovedApplicationsGetter interface {
	GetApprovedApplications() ([]models.Application, error)
}

func New(log *slog.Logger, approvedApplicationsGetter ApprovedApplicationsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.application.getApproved.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		applications, err := approvedApplicationsGetter.GetApprovedApplications()
		if err != nil {
			log.Error("failed to get approved applications", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get approved applications"))

			return
		}

		log.Info("get approved applications")

		render.JSON(w, r, applications)
		//responseOK(w, r, applications)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, applications []models.Application) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		Applications: applications,
	})
}

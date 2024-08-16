package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "projectsShowcase/internal/lib/api/response"
	"projectsShowcase/internal/lib/logger/sl"
)

type Request struct {
	ApplicantName           string `json:"applicant_name" validate:"required"`
	ApplicantEmail          string `json:"applicant_email" validate:"required"`
	ApplicantPhone          string `json:"applicant_phone" validate:"required"`
	PositionAndOrganization string `json:"position_and_organization" validate:"required"`
	ProjectDuration         string `json:"project_duration" validate:"required"`
	ProjectLevel            string `json:"project_level" validate:"required"`
	ProblemHolder           string `json:"problem_holder" validate:"required"`
	ProjectGoal             string `json:"project_goal" validate:"required"`
	Barrier                 string `json:"barrier" validate:"required"`
	ExistingSolutions       string `json:"existing_solutions" validate:"required"`
	Keywords                string `json:"keywords"`
	InterestedParties       string `json:"interested_parties" validate:"required"`
	Consultants             string `json:"consultants"`
	AdditionalMaterials     string `json:"additional_materials"`
	ProjectName             string `json:"project_name"`
}

type Response struct {
	resp.Response
	ID int64 `json:"id,omitempty"`
}

type ApplicationSaver interface {
	SaveApplication(applicantName, applicantEmail, applicantPhone, positionAndOrganization, projectDuration, projectLevel, problemHolder, projectGoal, barrier, existingSolutions, keywords, interestedParties, consultants, additionalMaterials, projectName, status string) (int64, error)
}

func New(log *slog.Logger, applicationSaver ApplicationSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.application.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
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

		id, err := applicationSaver.SaveApplication(req.ApplicantName, req.ApplicantEmail, req.ApplicantPhone, req.PositionAndOrganization, req.ProjectDuration, req.ProjectLevel, req.ProblemHolder, req.ProjectGoal, req.Barrier, req.ExistingSolutions, req.Keywords, req.InterestedParties, req.Consultants, req.AdditionalMaterials, req.ProjectName, "На рассмотрении")
		if err != nil {
			log.Error("failed to add application", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add application"))

			return
		}

		log.Info("application added", slog.Int64("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		ID:       id,
	})
}

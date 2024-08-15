package models

import "time"

type Application struct {
	ID                      int64
	ApplicantName           string
	ApplicantEmail          string
	ApplicantPhone          string
	PositionAndOrganization string
	ProjectDuration         string
	ProjectLevel            string
	ProblemHolder           string
	ProjectGoal             string
	Barrier                 string
	ExistingSolutions       string
	Keywords                string
	InterestedParties       string
	Consultants             string
	AdditionalMaterials     string
	ProjectName             string
	Status                  string
	SubmissionDate          time.Time
}

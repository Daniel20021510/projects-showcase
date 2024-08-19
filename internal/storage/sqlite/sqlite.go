package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"projectsShowcase/internal/domain/models"
	"projectsShowcase/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new SQLite storage instance.
//
// storagePath is the path to the SQLite database file.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		applicant_name TEXT NOT NULL,
		applicant_email TEXT NOT NULL,
		applicant_phone TEXT NOT NULL,
		position_and_organization TEXT NOT NULL,
		project_duration TEXT CHECK(project_duration IN ('1 семестр', '2 семестра')) NOT NULL,
		project_level TEXT CHECK(project_level IN ('Диагностический проект', 'Учебный проект', 'Учебно-прикладной проект', 'Прикладной проект')) NOT NULL,
		problem_holder TEXT NOT NULL,
		project_goal TEXT NOT NULL,
		barrier TEXT NOT NULL,
		existing_solutions TEXT NOT NULL,
		keywords TEXT,
		interested_parties TEXT,
		consultants TEXT,
		additional_materials TEXT,
		project_name TEXT NOT NULL,
		status TEXT CHECK(status IN ('На рассмотрении', 'Допущена', 'Удалена')) NOT NULL,
		submission_date DATETIME DEFAULT CURRENT_TIMESTAMP);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveApplication saves an application to the database.
//
// The function returns the ID of the inserted application (int64) and an error (error).
func (s *Storage) SaveApplication(
	applicantName,
	applicantEmail,
	applicantPhone,
	positionAndOrganization,
	projectDuration,
	projectLevel,
	problemHolder,
	projectGoal,
	barrier,
	existingSolutions,
	keywords,
	interestedParties,
	consultants,
	additionalMaterials,
	projectName,
	status string) (int64, error) {
	const op = "storage.sqlite.SaveApplication"

	stmt, err := s.db.Prepare(`INSERT INTO applications(
                         applicant_name,
                         applicant_email,
                         applicant_phone,
                         position_and_organization,
                         project_duration,
                         project_level,
                         problem_holder,
                         project_goal,
                         barrier,
                         existing_solutions,
                         keywords,
                         interested_parties,
                         consultants,
                         additional_materials,
                         project_name,
                         status)
					values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(
		applicantName,
		applicantEmail,
		applicantPhone,
		positionAndOrganization,
		projectDuration,
		projectLevel,
		problemHolder,
		projectGoal,
		barrier,
		existingSolutions,
		keywords,
		interestedParties,
		consultants,
		additionalMaterials,
		projectName,
		status,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to getApproved last insert id: %w", op, err)
	}

	return id, nil
}

// GetApplication retrieves an application from the database by its ID.
func (s *Storage) GetApplication(id int64) (models.Application, error) {
	const op = "storage.sqlite.GetApplication"

	stmt, err := s.db.Prepare(`SELECT 
    applicant_name,
	applicant_email,
	applicant_phone,
	position_and_organization,
	project_duration,
	project_level,
	problem_holder,
	project_goal,
	barrier,
	existing_solutions,
	keywords,
	interested_parties,
	consultants,
	additional_materials,
	project_name,
	status
    FROM applications WHERE id = ?`)
	if err != nil {
		return models.Application{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var application models.Application

	err = stmt.QueryRow(id).Scan(
		&application.ApplicantName,
		&application.ApplicantEmail,
		&application.ApplicantPhone,
		&application.PositionAndOrganization,
		&application.ProjectDuration,
		&application.ProjectLevel,
		&application.ProblemHolder,
		&application.ProjectGoal,
		&application.Barrier,
		&application.ExistingSolutions,
		&application.Keywords,
		&application.InterestedParties,
		&application.Consultants,
		&application.AdditionalMaterials,
		&application.ProjectName,
		&application.Status,
		&application.SubmissionDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Application{}, storage.ErrApplicationNotFound
	}
	if err != nil {
		return models.Application{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return application, nil
}

// GetApprovedApplications retrieves a list of approved applications from the database.
func (s *Storage) GetApprovedApplications() ([]models.Application, error) {
	const op = "storage.sqlite.GetApplication"

	stmt, err := s.db.Prepare(`SELECT 
	id,
    applicant_name,
	applicant_email,
	applicant_phone,
	position_and_organization,
	project_duration,
	project_level,
	problem_holder,
	project_goal,
	barrier,
	existing_solutions,
	keywords,
	interested_parties,
	consultants,
	additional_materials,
	project_name
    FROM applications WHERE status = ?
	ORDER BY submission_date`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var applications []models.Application

	rows, err := stmt.Query("Допущена")
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var application models.Application
		err = rows.Scan(
			&application.ID,
			&application.ApplicantName,
			&application.ApplicantEmail,
			&application.ApplicantPhone,
			&application.PositionAndOrganization,
			&application.ProjectDuration,
			&application.ProjectLevel,
			&application.ProblemHolder,
			&application.ProjectGoal,
			&application.Barrier,
			&application.ExistingSolutions,
			&application.Keywords,
			&application.InterestedParties,
			&application.Consultants,
			&application.AdditionalMaterials,
			&application.ProjectName,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		applications = append(applications, application)
	}

	return applications, nil
}

// DeleteApplication deletes the request from the database by its ID.
func (s *Storage) DeleteApplication(id int64) error {
	const op = "storage.sqlite.DeleteApplication"

	stmt, err := s.db.Prepare(`DELETE FROM applications WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to getApproved rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no application found with ID %d", id)
	}

	return nil
}

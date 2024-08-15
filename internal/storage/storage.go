package storage

import "errors"

var (
	ErrApplicationNotFound = errors.New("application not found")
	ErrProjectDuration     = errors.New("duration is not valid")
	ErrProjectLevel        = errors.New("level is not valid")
	ErrProjectStatus       = errors.New("status is not valid")
)

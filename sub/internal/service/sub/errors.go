package sub

import "errors"

var (
	ErrInvalidEmail  = errors.New("invalid subscriber's email")
	ErrAlreadyExists = errors.New("subscriber already exists")
	ErrFailedToSave  = errors.New("failed to save subsriber")
)

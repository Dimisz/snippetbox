package models

import "errors"

var (
	// ErrNoRecord is returned when no record was found
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials is returned if a user is trying to log in with an incorrect email address or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail is returned if a new user is trying to signup with an email address that's already in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

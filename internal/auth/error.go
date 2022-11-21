package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrLoginBusy          = errors.New("this login is busy")
	ErrUserNotLoggedIn    = errors.New("user not logged in")
)

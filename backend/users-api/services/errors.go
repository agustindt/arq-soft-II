package services

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrAccountDisabled       = errors.New("account disabled")
	ErrInvalidSecretKey      = errors.New("invalid secret key")
	ErrRootAlreadyExists     = errors.New("root user already exists")
	ErrCannotModifyRoot      = errors.New("cannot modify root user")
	ErrCannotDeactivateRoot  = errors.New("cannot deactivate root user")
	ErrCannotDeleteRoot      = errors.New("cannot delete root user")
	ErrAvatarUpdateFailed    = errors.New("failed to update avatar")
)

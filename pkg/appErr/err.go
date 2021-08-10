package appErr

import (
	"errors"
)

var (
	ErrUserIdNotFound = errors.New("user ID not found")
	ErrCardIdNotFound = errors.New("card ID not found")
	ErrDB             = errors.New("error db")
)

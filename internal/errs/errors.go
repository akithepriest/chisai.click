package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal = func(err error) error {
		return fmt.Errorf("intrnal error: %v", err)
	}
	ErrInvalid = errors.New("invalid params")
)
package constants

import "errors"

var (
	// ErrInternalServer is error for internal server error
	ErrInternalServer = errors.New("Internal Server Error")
	// ErrNotFound is error for data not found
	ErrNotFound = errors.New("Data Not Found")
)

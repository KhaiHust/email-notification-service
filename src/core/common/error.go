package common

import "net/http"

type Errs struct {
	Message    string
	Code       int
	StatusCode int
}

const (
	ErrEmailProviderNotFoundMessage = "Email provider not found"
)
const (
	ErrEmailProviderNotFoundCode = 404001
)

var (
	ErrEmailProviderNotFound = &Errs{
		Message:    ErrEmailProviderNotFoundMessage,
		Code:       ErrEmailProviderNotFoundCode,
		StatusCode: http.StatusNotFound,
	}
)

func (e *Errs) Error() string {
	return e.Message
}

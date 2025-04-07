package common

import "net/http"

var _ error = (*Errs)(nil)

type Errs struct {
	Message        string
	Code           int
	HttpStatusCode int
}

func New(message string, code int, httpStatusCode int) *Errs {
	return &Errs{
		Message:        message,
		Code:           code,
		HttpStatusCode: httpStatusCode,
	}
}

const (
	ErrEmailProviderNotFoundMessage   = "Email provider not found"
	ErrInternalServerMessage          = "Internal server error"
	EmailProviderParamNotFoundMessage = "Email provider not found"
	ErrRecordNotFoundMessage          = "Record not found"
	ErrForbiddenMessage               = "Forbidden"
	ErrOAuthCodeNotFoundMessage       = "OAuth code not found"
	ErrBadRequestMessage              = "Bad request"
)
const (
	ErrEmailProviderNotFoundCode   = 404001
	ErrInternalServerCode          = 500000
	ErrBadRequestCode              = 400000
	EmailProviderParamNotFoundCode = 400001
	ErrOAuthCodeNotFoundCode       = 400002
	ErrRecordNotFoundCode          = 404002
	ErrForbiddenCode               = 403001
)

var (
	ErrEmailProviderNotFound = Errs{
		Message:        ErrEmailProviderNotFoundMessage,
		Code:           ErrEmailProviderNotFoundCode,
		HttpStatusCode: http.StatusNotFound,
	}
	ErrInternalServer = Errs{
		Message:        ErrInternalServerMessage,
		Code:           ErrInternalServerCode,
		HttpStatusCode: http.StatusInternalServerError,
	}
	ErrEmailProviderParamNotFound = Errs{
		Message:        EmailProviderParamNotFoundMessage,
		Code:           EmailProviderParamNotFoundCode,
		HttpStatusCode: http.StatusBadRequest,
	}
	ErrRecordNotFound = Errs{
		Message:        ErrRecordNotFoundMessage,
		Code:           ErrRecordNotFoundCode,
		HttpStatusCode: http.StatusNotFound,
	}
	ErrForbidden = Errs{
		Message:        ErrForbiddenMessage,
		Code:           ErrForbiddenCode,
		HttpStatusCode: http.StatusForbidden,
	}
	ErrOAuthCodeNotFound = Errs{
		Message:        ErrOAuthCodeNotFoundMessage,
		Code:           ErrOAuthCodeNotFoundCode,
		HttpStatusCode: http.StatusBadRequest,
	}
	ErrBadRequest = Errs{
		Message:        ErrBadRequestMessage,
		Code:           ErrBadRequestCode,
		HttpStatusCode: http.StatusBadRequest,
	}
)

func (e Errs) Error() string {
	return e.Message
}

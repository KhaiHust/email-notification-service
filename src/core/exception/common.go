package exception

import (
	"github.com/golibs-starter/golib/exception"
	"net/http"
)

var (
	InternalServerException = exception.New(http.StatusInternalServerError, "Internal server error")
)

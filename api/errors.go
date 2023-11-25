package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorWithData struct {
	error
	Data any
}

func (e ErrorWithData) Error() string {
	return e.Error()
}

var (
	ErrBadRequest      = errors.New("base.bad_request")
	ErrValidation      = errors.New("base.validation")
	ErrUnauthorized    = errors.New("base.unauthorized")
	ErrForbidden       = errors.New("base.forbidden")
	ErrNotFound        = errors.New("base.not_found")
	ErrConflict        = errors.New("base.conflict")
	ErrTooManyRequests = errors.New("base.too_many_request")
	ErrInternal        = errors.New("base.internal")
)

func BaseError(base error, text string) error {
	return fmt.Errorf("%w: %s", base, text)
}

func ErrorData(err error, data any) error {
	return ErrorWithData{
		error: err,
		Data:  data,
	}
}

func ServeError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var errorWithData ErrorWithData
	if errors.Is(err, &errorWithData) {
		RespondErrorStatusData(c, ErrorStatusCode(err), ErrorResponseData(errorWithData.Data, err), err)
		return
	}

	RespondError(c, err)
}

func ErrorStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrTooManyRequests):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

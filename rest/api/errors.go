package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/errx"
)

var (
	ErrBadRequest          = errors.New("base.bad_request")
	ErrValidation          = errors.New("base.validation")
	ErrUnauthorized        = errors.New("base.unauthorized")
	ErrForbidden           = errors.New("base.forbidden")
	ErrNotFound            = errors.New("base.not_found")
	ErrConflict            = errors.New("base.conflict")
	ErrTooManyRequests     = errors.New("base.too_many_request")
	ErrInternalServerError = errors.New("base.internal_server_error")
)

func ServeError(c *gin.Context, err error) {
	if err == nil {
		RespondOK(c)
		return
	}

	var errorWithData errx.Error
	if errors.As(err, &errorWithData) && errorWithData.ResponseData != nil {
		RespondErrorData(c, err, errorWithData.ResponseData)
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
	case errors.Is(err, ErrInternalServerError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

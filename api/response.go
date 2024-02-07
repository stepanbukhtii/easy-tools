package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrorCodeInvalidRoute   = "invalid_route"
	ErrorCodeInvalidRequest = "invalid_request"
	ErrorCodeInvalidQuery   = "invalid_query"
	ErrorCodeInvalidBody    = "invalid_body"
	ErrorCodeInvalidServer  = "invalid_server"
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
	Pages int64  `json:"pages,omitempty"`
}

type Error struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

func RespondOK(c *gin.Context) {
	c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func RespondData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Data: data})
}

func RespondDataPages(c *gin.Context, data any, pages int64) {
	c.JSON(http.StatusOK, Response{
		Data:  data,
		Pages: pages,
	})
}

func ErrorResponse(code string, err error) Response {
	var description string
	if err != nil {
		description = err.Error()
	}
	return Response{
		Error: &Error{
			Code:        code,
			Description: description,
		},
	}
}

func RespondBadRequest(c *gin.Context, code string, desc error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(code, desc))
}

func RespondUnauthorized(c *gin.Context, code string, desc error) {
	c.Header(HeaderAuthenticate, "Bearer realm=\"api\"")
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(code, desc))
}

func RespondForbidden(c *gin.Context, code string, desc error) {
	c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse(code, desc))
}

func RespondNotFound(c *gin.Context, code string, desc error) {
	c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(code, desc))
}

func RespondConflict(c *gin.Context, code string, desc error) {
	c.AbortWithStatusJSON(http.StatusConflict, ErrorResponse(code, desc))
}

func RespondInternalError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

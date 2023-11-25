package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stepanbukhtii/easy-tools/easycontext"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Meta   any    `json:"meta,omitempty"`
	Error  error  `json:"error,omitempty"`
}

func ErrorResponse(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}

func ErrorResponseData(data any, err error) Response {
	return Response{
		Status: StatusError,
		Data:   data,
		Error:  err,
	}
}

type MetaPages struct {
	Pages int64 `json:"pages,omitempty"`
}

func RespondOK(c *gin.Context) {
	RespondJSON(c, http.StatusOK, Response{Status: StatusSuccess})
}

func RespondData(c *gin.Context, data any) {
	RespondJSON(c, http.StatusOK, Response{
		Status: StatusSuccess,
		Data:   data,
	})
}

func RespondDataPages(c *gin.Context, data any, pages int64) {
	RespondJSON(c, http.StatusOK, Response{
		Status: StatusSuccess,
		Data:   data,
		Meta: MetaPages{
			Pages: pages,
		},
	})
}

func RespondDataMeta(c *gin.Context, data, meta any) {
	RespondJSON(c, http.StatusOK, Response{
		Status: StatusSuccess,
		Data:   data,
		Meta:   meta,
	})
}

func RespondBadRequest(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusBadRequest, err)
}

func RespondUnauthorized(c *gin.Context, err error) {
	c.Header(HeaderAuthenticate, "Bearer realm=\"api\"")
	RespondStatusError(c, http.StatusUnauthorized, err)
}

func RespondForbidden(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusForbidden, err)
}

func RespondNotFound(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusNotFound, err)
}

func RespondConflict(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusConflict, err)
}

func RespondToManyRequests(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusTooManyRequests, err)
}

func RespondInternalError(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusInternalServerError, err)
}

func RespondJSON(c *gin.Context, code int, data any) {
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, easycontext.Locale(c.Request.Context()))
	c.JSON(code, data)
}

func RespondError(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, easycontext.Locale(c.Request.Context()))
	c.AbortWithStatusJSON(ErrorStatusCode(err), ErrorResponse(err))
}

func RespondStatusError(c *gin.Context, code int, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, easycontext.Locale(c.Request.Context()))
	c.AbortWithStatusJSON(code, ErrorResponse(err))
}

func RespondErrorStatusData(c *gin.Context, code int, data any, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, easycontext.Locale(c.Request.Context()))
	c.AbortWithStatusJSON(code, data)
}

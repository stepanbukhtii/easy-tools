package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stepanbukhtii/easy-tools/econtext"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
	Meta   any    `json:"meta,omitempty"`
}

type MetaPages struct {
	Pages int64 `json:"pages,omitempty"`
}

func NewResponse(status string, data any) Response {
	return Response{
		Status: status,
		Data:   data,
	}
}

func NewResponseMeta(status string, data, meta any) Response {
	return Response{
		Status: status,
		Data:   data,
		Meta:   meta,
	}
}

func NewErrorResponse(err error) Response {
	return Response{
		Status: StatusError,
		Error:  errorText(err),
	}
}

func NewErrorResponseData(err error, data any) Response {
	return Response{
		Status: StatusError,
		Error:  errorText(err),
		Data:   data,
	}
}

func RespondOK(c *gin.Context) {
	RespondJSON(c, http.StatusOK, NewResponse(StatusSuccess, nil))
}

func RespondData(c *gin.Context, data any) {
	RespondJSON(c, http.StatusOK, NewResponse(StatusSuccess, data))
}

func RespondDataPages(c *gin.Context, data any, pages int64) {
	RespondJSON(c, http.StatusOK, NewResponseMeta(StatusSuccess, data, MetaPages{Pages: pages}))
}

func RespondDataMeta(c *gin.Context, data, meta any) {
	RespondJSON(c, http.StatusOK, NewResponseMeta(StatusSuccess, data, meta))
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

func RespondTooManyRequests(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusTooManyRequests, err)
}

func RespondInternalError(c *gin.Context, err error) {
	RespondStatusError(c, http.StatusInternalServerError, err)
}

func RespondJSON(c *gin.Context, code int, data any) {
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, econtext.ClientInfo(c.Request.Context()).Locale)
	c.JSON(code, data)
}

func RespondError(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, econtext.ClientInfo(c.Request.Context()).Locale)
	c.AbortWithStatusJSON(ErrorStatusCode(err), NewErrorResponse(err))
}

func RespondStatusError(c *gin.Context, code int, err error) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, econtext.ClientInfo(c.Request.Context()).Locale)
	c.AbortWithStatusJSON(code, NewErrorResponse(err))
}

func RespondErrorData(c *gin.Context, err error, data any) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, econtext.ClientInfo(c.Request.Context()).Locale)
	c.AbortWithStatusJSON(ErrorStatusCode(err), NewErrorResponseData(err, data))
}

func RespondStatusErrorData(c *gin.Context, code int, err error, data any) {
	if err != nil {
		_ = c.Error(err)
	}
	c.Header(HeaderContentTypeOption, "nosniff")
	c.Header(HeaderContentLanguage, econtext.ClientInfo(c.Request.Context()).Locale)
	c.AbortWithStatusJSON(code, NewErrorResponseData(err, data))
}

// for remove base error text use return regexp.MustCompile(`:\sbase\..+$`).ReplaceAllString(err.Error(), "")
func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

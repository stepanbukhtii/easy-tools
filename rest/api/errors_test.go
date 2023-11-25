package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/errx"
	"github.com/stretchr/testify/require"
)

func TestServeError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testResponseData := map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "without error",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		}, {
			name:           "bad request",
			err:            ErrBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","error":"base.bad_request"}`,
		}, {
			name:           "validation",
			err:            ErrValidation,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","error":"base.validation"}`,
		}, {
			name:           "unauthorized",
			err:            ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":"error","error":"base.unauthorized"}`,
		}, {
			name:           "forbidden",
			err:            ErrForbidden,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"status":"error","error":"base.forbidden"}`,
		}, {
			name:           "not found",
			err:            ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":"error","error":"base.not_found"}`,
		}, {
			name:           "conflict",
			err:            ErrConflict,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":"error","error":"base.conflict"}`,
		}, {
			name:           "too many request",
			err:            ErrTooManyRequests,
			expectedStatus: http.StatusTooManyRequests,
			expectedBody:   `{"status":"error","error":"base.too_many_request"}`,
		}, {
			name:           "internal server error",
			err:            ErrInternalServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"error","error":"base.internal_server_error"}`,
		}, {
			name:           "internal server error unknow error",
			err:            errors.New("test error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"error","error":"test error"}`,
		}, {
			name:           "errx package error",
			err:            errx.Wrap(ErrBadRequest, "custom_error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","error":"custom_error: base.bad_request"}`,
		}, {
			name:           "errx package error with wrapped error",
			err:            errx.Wrap(errx.Wrap(ErrBadRequest, "custom_error"), "wrap_error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","error":"wrap_error: custom_error: base.bad_request"}`,
		}, {
			name:           "errx package error with response data",
			err:            errx.Wrap(ErrBadRequest, "custom_error").WithResponseData(testResponseData),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","error":"custom_error: base.bad_request","data":{"key1":"value1","key2":"value2","key3":"value3"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, gin.New())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

			ServeError(c, test.err)

			require.Equal(t, test.expectedStatus, c.Writer.Status())
			require.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

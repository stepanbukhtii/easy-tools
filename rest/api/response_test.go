package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	responseData := map[string]string{"key1": "value1", "key2": "value2"}
	metaData := MetaPages{Pages: 3}
	responseErr := errors.New("error text")

	tests := []struct {
		name           string
		responseFunc   func(c *gin.Context)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "empty response",
			responseFunc:   func(c *gin.Context) { RespondOK(c) },
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		}, {
			name:           "with data",
			responseFunc:   func(c *gin.Context) { RespondData(c, responseData) },
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","data":{"key1":"value1","key2":"value2"}}`,
		}, {
			name:           "with data and pages",
			responseFunc:   func(c *gin.Context) { RespondDataPages(c, responseData, 2) },
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","data":{"key1":"value1","key2":"value2"},"meta":{"pages":2}}`,
		}, {
			name:           "with data and meta",
			responseFunc:   func(c *gin.Context) { RespondDataMeta(c, responseData, metaData) },
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success","data":{"key1":"value1","key2":"value2"},"meta":{"pages":3}}`,
		}, {
			name:           "error bad request",
			responseFunc:   func(c *gin.Context) { RespondBadRequest(c, responseErr) },
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name:           "error unauthorized",
			responseFunc:   func(c *gin.Context) { RespondUnauthorized(c, responseErr) },
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name:           "error forbidden",
			responseFunc:   func(c *gin.Context) { RespondForbidden(c, responseErr) },
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name:           "error not found",
			responseFunc:   func(c *gin.Context) { RespondNotFound(c, responseErr) },
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name:           "error conflict",
			responseFunc:   func(c *gin.Context) { RespondConflict(c, responseErr) },
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name:           "error too many requests",
			responseFunc:   func(c *gin.Context) { RespondTooManyRequests(c, responseErr) },
			expectedStatus: http.StatusTooManyRequests,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name:           "error internal server error",
			responseFunc:   func(c *gin.Context) { RespondInternalError(c, responseErr) },
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":"error","error":"error text"}`,
		}, {
			name: "error with status, error and data",
			responseFunc: func(c *gin.Context) {
				RespondStatusErrorData(c, http.StatusConflict, responseErr, responseData)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"status":"error","error":"error text","data":{"key1":"value1","key2":"value2"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, gin.New())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

			test.responseFunc(c)

			require.Equal(t, test.expectedStatus, c.Writer.Status())
			require.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

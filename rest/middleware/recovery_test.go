package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tests := []struct {
		name           string
		handler        gin.HandlerFunc
		expectedStatus int
	}{
		{
			name:           "success",
			handler:        func(c *gin.Context) { api.RespondOK(c) },
			expectedStatus: http.StatusOK,
		}, {
			name:           "panic",
			handler:        func(c *gin.Context) { panic("panic text") },
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			_, router := gin.CreateTestContext(recorder)

			router.GET("/", Recovery, test.handler)

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			assert.NoError(t, err)

			router.ServeHTTP(recorder, req)

			require.Equal(t, test.expectedStatus, recorder.Code)
			fmt.Println(recorder.Body.String())
		})
	}
}

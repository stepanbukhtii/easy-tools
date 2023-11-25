package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tests := []struct {
		name string
	}{
		{
			"success",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte{}))

			Logger(c)
		})
	}
}

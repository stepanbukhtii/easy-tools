package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/elog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	slog.SetDefault(slog.New(elog.NewCustomHandler(slog.NewJSONHandler(os.Stdout, nil), nil)))
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tests := []struct {
		name        string
		traceParent string
	}{
		{
			name:        "success",
			traceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", io.NopCloser(bytes.NewReader(nil)))
			c.Request.Header.Set("traceparent", test.traceParent)

			Trace(c)
			Logger(c)
		})
	}
}

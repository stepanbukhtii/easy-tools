package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestTrace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tests := []struct {
		name            string
		traceParent     string
		expectedTraceID string
		expectedSpanID  string
	}{
		{
			name:            "success",
			traceParent:     "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			expectedTraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
			expectedSpanID:  "00f067aa0ba902b7",
		}, {
			name:            "not found",
			traceParent:     "",
			expectedTraceID: "00000000000000000000000000000000",
			expectedSpanID:  "0000000000000000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			if test.traceParent != "" {
				c.Request.Header.Set("traceparent", test.traceParent)
			}

			Trace(c)

			ctx := c.Request.Context()

			require.Equal(t, test.expectedTraceID, trace.SpanContextFromContext(ctx).TraceID().String())
			require.Equal(t, test.expectedSpanID, trace.SpanContextFromContext(ctx).SpanID().String())
		})
	}
}

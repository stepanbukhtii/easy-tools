package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Trace(c *gin.Context) {
	ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

	if trace.SpanContextFromContext(ctx).IsValid() {
		c.Request = c.Request.WithContext(ctx)
	}
}

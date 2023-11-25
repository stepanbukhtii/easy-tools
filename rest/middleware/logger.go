package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func Logger(c *gin.Context) {
	ctx := c.Request.Context()

	if econtext.SkipLogger(ctx) {
		c.Next()
		return
	}

	start := time.Now().UTC()

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.With(elog.Err(err)).ErrorContext(ctx, "read request body failed")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	c.Next()

	duration := time.Now().UTC().Sub(start)

	logger := econtext.Logger(ctx).With(
		slog.String(string(semconv.HTTPRequestMethodKey), c.Request.Method),
		slog.String(string(semconv.HTTPResponseStatusCodeKey), strconv.Itoa(c.Writer.Status())),
		slog.Duration(elog.HTTPRequestDuration, duration),
		slog.String(string(semconv.URLPathKey), c.Request.URL.RequestURI()),
		slog.String(string(semconv.ClientAddressKey), c.ClientIP()),
		slog.String(string(semconv.UserAgentOriginalKey), c.Request.UserAgent()),
	)

	if len(bodyBytes) > 0 {
		logger = logger.With(slog.String(elog.HTTPRequestBodyContent, string(bodyBytes)))
	}

	var loggerErr error
	if len(c.Errors) > 0 {
		loggerErr = fmt.Errorf("%s", c.Errors.String())
	}

	switch {
	case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
		logger.With(elog.Err(loggerErr)).WarnContext(ctx, "http request")
	case c.Writer.Status() >= http.StatusInternalServerError:
		logger.With(elog.Err(loggerErr)).ErrorContext(ctx, "http request")
	default:
		logger.InfoContext(ctx, "http request")
	}
}

func SkipLogger(c *gin.Context) {
	c.Request = c.Request.WithContext(econtext.SetSkipLogger(c.Request.Context(), true))
}

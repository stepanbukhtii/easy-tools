package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
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
		slog.With(elog.Err(err)).Error("failed to read request body in middleware logger")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	c.Next()

	latency := time.Now().UTC().Sub(start)

	logger := econtext.Logger(ctx)
	if !logger.Enabled(ctx, slog.LevelError) {
		logger = slog.Default()
	}

	logger = logger.With(
		slog.String(elog.HTTPRequestMethod, c.Request.Method),
		slog.Any(elog.HTTPRequestBodyContent, bodyBytes),
		slog.Int(elog.HTTPResponseStatusCode, c.Writer.Status()),
		slog.String(elog.URLPath, c.Request.URL.RequestURI()),
		slog.String(elog.ClientIP, c.ClientIP()),
		slog.String(elog.UserAgentOriginal, c.Request.UserAgent()),
		slog.Duration(elog.EventDuration, latency),
	)

	if traceID := econtext.TraceID(ctx); traceID != "" {
		logger = logger.With(slog.String(elog.TraceID, traceID))
	}

	var loggerErr error
	if len(c.Errors) > 0 {
		loggerErr = fmt.Errorf("%s", c.Errors.String())
	}

	switch {
	case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
		logger.With(elog.Err(loggerErr)).Warn("Request")
	case c.Writer.Status() >= http.StatusInternalServerError:
		logger.With(elog.Err(loggerErr)).Error("Request")
	default:
		logger.Info("Request")
	}
}

func SkipLogger(c *gin.Context) {
	c.Request = c.Request.WithContext(econtext.SetSkipLogger(c.Request.Context(), true))
}

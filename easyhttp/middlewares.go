package easyhttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/stepanbukhtii/easy-tools/api"
	"github.com/stepanbukhtii/easy-tools/easycontext"
	"github.com/stepanbukhtii/easy-tools/easylog"
)

const LocaleEN = "en"

var DefaultLocale = LocaleEN

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			var brokenPipe bool
			var err error
			switch x := r.(type) {
			case *net.OpError:
				// Check for a broken connection, as it is not really a condition that warrants a panic stack trace.
				var se *os.SyscallError
				if errors.As(x, &se) {
					seStr := strings.ToLower(se.Error())
					if strings.Contains(seStr, "broken pipe") ||
						strings.Contains(seStr, "connection reset by peer") {
						brokenPipe = true
					}
				}
				err = x.Err
			case string:
				err = errors.New(x)
			case error:
				err = x
			}

			log.Error().Err(err).Bytes(easylog.ErrorStackTrace, debug.Stack()).Msg("panic recovered")

			if brokenPipe {
				c.Abort()
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}()

	c.Next()
}

func MiddlewareLogger(c *gin.Context) {
	start := time.Now().UTC()

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	c.Next()

	ctx := c.Request.Context()

	if easycontext.SkipLogger(ctx) {
		return
	}

	latency := time.Now().UTC().Sub(start)

	logWith := log.With().
		Int(easylog.HTTPResponseStatusCode, c.Writer.Status()).
		Str(easylog.HTTPRequestMethod, c.Request.Method).
		Str(easylog.URLPath, c.Request.URL.RequestURI()).
		Str(easylog.ClientIP, c.ClientIP()).
		Dur(easylog.EventDuration, latency).
		Str(easylog.UserAgentOriginal, c.Request.UserAgent()).
		Bytes(easylog.HTTPRequestBodyContent, bodyBytes)

	if traceID := easycontext.TraceID(ctx); traceID != "" {
		logWith = logWith.Str(easylog.TraceID, traceID)
	}

	requestLogger := logWith.Logger()

	switch {
	case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
		requestLogger.Warn().Msg("Request")
	case c.Writer.Status() >= http.StatusInternalServerError:
		if len(c.Errors) > 0 {
			requestLogger.Err(fmt.Errorf("%s", c.Errors.String())).Msg("Request")
			return
		}

		requestLogger.Error().Msg("Request")
	default:
		requestLogger.Info().Msg("Request")
	}
}

func SkipLogger(c *gin.Context) {
	c.Request = c.Request.WithContext(easycontext.SetSkipLogger(c.Request.Context(), true))
	c.Next()
}

func TraceID(c *gin.Context) {
	traceID := c.GetHeader(api.HeaderTraceID)
	if traceID == "" {
		traceID = uuid.New().String()
		c.Header(api.HeaderTraceID, traceID)
	}

	c.Request = c.Request.WithContext(easycontext.SetTraceID(c.Request.Context(), traceID))
}

func Locale(c *gin.Context) {
	locale := c.GetHeader(api.HeaderAcceptLanguage)
	if locale == "" {
		locale = DefaultLocale
	}
	c.Request = c.Request.WithContext(easycontext.SetLocale(c.Request.Context(), locale))
}

func IPAddress(c *gin.Context) {
	c.Request = c.Request.WithContext(easycontext.SetIPAddress(c.Request.Context(), c.ClientIP()))
}

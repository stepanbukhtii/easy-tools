package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/elog"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			var brokenPipe bool
			var err error
			switch x := r.(type) {
			case *net.OpError:
				// Check for a broken connection, as it is not really a condition that warrants a panic stack trace.
				var se *os.SyscallError
				if errors.Is(x.Err, syscall.EPIPE) || errors.Is(x.Err, syscall.ECONNRESET) {
					brokenPipe = true
				}
				if errors.As(x, &se) {
					seStr := strings.ToLower(se.Error())
					if slices.Contains([]string{"broken pipe", "connection reset by peer"}, seStr) {
						brokenPipe = true
					}
				}
				err = x.Err
			case string:
				err = errors.New(x)
			case error:
				err = x
			}

			slog.With(
				slog.Any(elog.ErrorMessage, err),
				slog.Any(elog.ErrorStackTrace, debug.Stack()),
			).Error("panic recovered")

			if brokenPipe {
				c.Abort()
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		}
	}()

	c.Next()
}

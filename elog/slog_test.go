package elog

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/stepanbukhtii/easy-tools/errx"
	"github.com/stretchr/testify/require"
)

func TestSlog(t *testing.T) {
	logData := map[string]string{"key1": "value1", "key2": "value2"}

	var tests = []struct {
		name           string
		err            error
		withErrorTrace bool
		expectedLog    string
	}{
		{
			name:        "simple error",
			err:         errors.New("error text"),
			expectedLog: `{"level":"ERROR","message":"Message","error.type":"error text"}`,
		}, {
			name:        "simple custom error",
			err:         errx.New("error text"),
			expectedLog: `{"level":"ERROR","message":"Message","error.type":"error text"}`,
		}, {
			name:        "error with log data",
			err:         errx.New("error text").WithLogData(logData),
			expectedLog: `{"level":"ERROR","message":"Message","error":{"type":"error text","data":{"key1":"value1","key2":"value2"}}}`,
		}, {
			name: "error with log data map",
			err: errx.New("error text").
				AddLogStr("key3", "value3").
				AddLogStr("key4", "value4"),
			expectedLog: `{"level":"ERROR","message":"Message","error":{"type":"error text","data":{"key3":"value3","key4":"value4"}}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var logWriter bytes.Buffer

			opts := &slog.HandlerOptions{
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return slogReplaceAttr(groups, a)
				},
			}
			logger := slog.New(&ContextHandler{Handler: slog.NewJSONHandler(&logWriter, opts)})

			logger.With(Err(test.err)).Error("Message")

			require.Equal(t, test.expectedLog, strings.TrimSpace(logWriter.String()))
		})
	}
}

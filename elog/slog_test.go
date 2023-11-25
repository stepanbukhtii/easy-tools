package elog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stepanbukhtii/easy-tools/errx"
	"github.com/stretchr/testify/require"
)

func TestSlog(t *testing.T) {
	logData := map[string]string{"key1": "value1", "key2": "value2"}

	var tests = []struct {
		name        string
		err         error
		expectedLog string
	}{
		{
			name:        "simple error",
			err:         errx.New("error text"),
			expectedLog: `{"level":"ERROR","message":"Message","error.message":"error text"}`,
		}, {
			name:        "error with log data",
			err:         errx.New("error text").WithLogData(logData),
			expectedLog: `{"level":"ERROR","message":"Message","error":{"message":"error text","data":{"key1":"value1","key2":"value2"}}}`,
		}, {
			name: "error with log data map",
			err: errx.New("error text").
				AddLogStr("key3", "value3").
				AddLogStr("key4", "value4"),
			expectedLog: `{"level":"ERROR","message":"Message","error":{"message":"error text","data":{"key3":"value3","key4":"value4"}}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var logWriter bytes.Buffer

			opts := &slog.HandlerOptions{
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey || a.Key == ECSVersion {
						return slog.Attr{}
					}
					return slogReplaceAttr(groups, a)
				},
			}
			logger := slog.New(&ContextHandler{Handler: slog.NewJSONHandler(&logWriter, opts)})

			logger.Error("Message", slog.Any(ErrorMessage, test.err))

			require.Equal(t, test.expectedLog, strings.TrimSpace(logWriter.String()))
		})
	}
}

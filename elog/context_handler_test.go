package elog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/rest/client"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestContextHandler_Handle(t *testing.T) {
	var logWriter bytes.Buffer
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	logger := slog.New(NewCustomHandler(slog.NewJSONHandler(&logWriter, opts), nil))

	ctx := context.Background()
	clientInfo := client.Info{
		Subject:   "Subject",
		Roles:     []string{"Role"},
		Locale:    "Locale",
		IPAddress: "IPAddress",
		UserAgent: "UserAgent",
		DeviceID:  "DeviceID",
	}
	ctx = econtext.SetClientInfo(ctx, clientInfo)

	traceID, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	require.NoError(t, err)

	spanID, err := trace.SpanIDFromHex("00f067aa0ba902b7")
	require.NoError(t, err)

	spanContextConfig := trace.SpanContextConfig{TraceID: traceID, SpanID: spanID}

	ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(spanContextConfig))

	logger.With(slog.String("key", "value")).InfoContext(ctx, "message")

	expectedLog := `{"level":"INFO","msg":"message","key":"value","user.id":"Subject","user.roles":["Role"],"client.address":"IPAddress","user_agent.original":"UserAgent","device.id":"DeviceID","trace.id":"4bf92f3577b34da6a3ce929d0e0e4736","span.id":"00f067aa0ba902b7"}`
	require.Equal(t, expectedLog, strings.TrimSpace(logWriter.String()))
}

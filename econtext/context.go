package econtext

import (
	"context"
	"log/slog"

	"github.com/stepanbukhtii/easy-tools/rest/client"
	"go.opentelemetry.io/otel/trace"
)

type (
	contextKey[T any] struct{}
	loggerKey         struct{}
	clientInfoKey     struct{}
	skipLoggerKey     struct{}
	logFieldsKey      struct{}
)

func WithValue[T any](ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, contextKey[T]{}, value)
}

func Value[T1, T2 any](ctx context.Context) T2 {
	var value T2
	value, _ = ctx.Value(contextKey[T1]{}).(T2)
	return value
}

func ValueDefault[T1, T2 any](ctx context.Context, defaultValue T2) T2 {
	var value T2
	var ok bool

	value, ok = ctx.Value(contextKey[T1]{}).(T2)
	if !ok {
		return defaultValue
	}

	return value
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		if l == logger {
			return ctx
		}
	}
	return context.WithValue(ctx, loggerKey{}, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.New(slog.DiscardHandler)
}

func TraceID(ctx context.Context) string {
	traceSpan := trace.SpanContextFromContext(ctx)
	if !traceSpan.IsValid() || !traceSpan.HasTraceID() {
		return ""
	}

	return traceSpan.TraceID().String()
}

func SpanID(ctx context.Context) string {
	traceSpan := trace.SpanContextFromContext(ctx)
	if !traceSpan.IsValid() || !traceSpan.HasSpanID() {
		return ""
	}

	return traceSpan.SpanID().String()
}

func SetClientInfo(ctx context.Context, clientInfo client.Info) context.Context {
	return context.WithValue(ctx, clientInfoKey{}, clientInfo)
}

func ClientInfo(ctx context.Context) client.Info {
	clientInfo, ok := ctx.Value(clientInfoKey{}).(client.Info)
	if !ok {
		return client.Info{}
	}
	return clientInfo
}

func SetSubject(ctx context.Context, subject string) context.Context {
	clientInfo, ok := ctx.Value(clientInfoKey{}).(client.Info)
	if ok {
		clientInfo.Subject = subject
		return SetClientInfo(ctx, clientInfo)
	}
	return SetClientInfo(ctx, client.Info{Subject: subject})
}

func SetRoles(ctx context.Context, roles []string) context.Context {
	clientInfo, ok := ctx.Value(clientInfoKey{}).(client.Info)
	if ok {
		clientInfo.Roles = roles
		return SetClientInfo(ctx, clientInfo)
	}
	return SetClientInfo(ctx, client.Info{Roles: roles})
}

func SetSkipLogger(ctx context.Context, skipLogger bool) context.Context {
	return context.WithValue(ctx, skipLoggerKey{}, skipLogger)
}

func SkipLogger(ctx context.Context) bool {
	skipLogger, ok := ctx.Value(skipLoggerKey{}).(bool)
	if !ok {
		return false
	}
	return skipLogger
}

func SetLogField(ctx context.Context, key, value string) context.Context {
	if fields, ok := ctx.Value(logFieldsKey{}).(map[string]string); ok {
		fields[key] = value
		return context.WithValue(ctx, logFieldsKey{}, fields)
	}
	return context.WithValue(ctx, logFieldsKey{}, map[string]string{key: value})
}

func LogFields(ctx context.Context) map[string]string {
	fields, ok := ctx.Value(logFieldsKey{}).(map[string]string)
	if !ok {
		return nil
	}
	return fields
}

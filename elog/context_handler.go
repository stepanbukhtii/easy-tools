package elog

import (
	"context"
	"log/slog"

	"github.com/stepanbukhtii/easy-tools/econtext"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

type CustomHandler func(context.Context, slog.Record) error

type ContextHandler struct {
	slog.Handler
	CustomHandler CustomHandler
}

func NewCustomHandler(baseHandler slog.Handler, customHandler CustomHandler) *ContextHandler {
	return &ContextHandler{
		Handler:       baseHandler,
		CustomHandler: customHandler,
	}
}

func (h *ContextHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	clientInfo := econtext.ClientInfo(ctx)
	if clientInfo.Subject != "" {
		r.AddAttrs(slog.String(string(semconv.UserIDKey), clientInfo.Subject))
	}
	if clientInfo.Roles != nil {
		r.AddAttrs(slog.Any(string(semconv.UserRolesKey), clientInfo.Roles))
	}
	if clientInfo.IPAddress != "" {
		r.AddAttrs(slog.String(string(semconv.ClientAddressKey), clientInfo.IPAddress))
	}
	if clientInfo.UserAgent != "" {
		r.AddAttrs(slog.String(string(semconv.UserAgentOriginalKey), clientInfo.UserAgent))
	}
	if clientInfo.DeviceID != "" {
		r.AddAttrs(slog.String(string(semconv.DeviceIDKey), clientInfo.DeviceID))
	}

	if traceSpan := trace.SpanContextFromContext(ctx); traceSpan.IsValid() {
		if traceSpan.HasTraceID() {
			r.AddAttrs(slog.String("trace.id", traceSpan.TraceID().String()))
		}
		if traceSpan.HasSpanID() {
			r.AddAttrs(slog.String("span.id", traceSpan.SpanID().String()))
		}
	}

	if fields := econtext.LogFields(ctx); fields != nil {
		for key, value := range fields {
			r.AddAttrs(slog.String(key, value))
		}
	}

	if h.CustomHandler != nil {
		if err := h.CustomHandler(ctx, r); err != nil {
			return err
		}
	}

	return h.Handler.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{
		Handler:       h.Handler.WithAttrs(attrs),
		CustomHandler: h.CustomHandler,
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		Handler:       h.Handler.WithGroup(name),
		CustomHandler: h.CustomHandler,
	}
}

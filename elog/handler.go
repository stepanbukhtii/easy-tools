package elog

import (
	"context"
	"log/slog"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"go.opentelemetry.io/otel/trace"
)

type CustomHandle func(context.Context, slog.Record) error

type ContextHandler struct {
	slog.Handler
	CustomHandle func(context.Context, slog.Record) error
}

func NewCustomHandler(baseHandler slog.Handler, customHandle CustomHandle) *ContextHandler {
	return &ContextHandler{
		Handler:      baseHandler,
		CustomHandle: customHandle,
	}
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceSpan := trace.SpanContextFromContext(ctx); traceSpan.IsValid() {
		if traceSpan.HasTraceID() {
			r.AddAttrs(slog.String(TraceID, traceSpan.TraceID().String()))
		}
		if traceSpan.HasSpanID() {
			r.AddAttrs(slog.String(SpanID, traceSpan.SpanID().String()))
		}
	}

	clientInfo := econtext.ClientInfo(ctx)
	if clientInfo.Subject != "" {
		r.AddAttrs(slog.String(UserID, clientInfo.Subject))
	}
	if clientInfo.Roles != nil {
		r.AddAttrs(slog.Any(UserRoles, clientInfo.Roles))
	}
	if clientInfo.IPAddress != "" {
		r.AddAttrs(slog.String(ClientIP, clientInfo.IPAddress))
	}
	if clientInfo.UserAgent != "" {
		r.AddAttrs(slog.String(UserAgentName, clientInfo.UserAgent))
	}
	if clientInfo.DeviceID != "" {
		r.AddAttrs(slog.String(DeviceID, clientInfo.DeviceID))
	}

	if fields := econtext.LogFields(ctx); fields != nil {
		for key, value := range fields {
			r.AddAttrs(slog.String(key, value))
		}
	}

	if h.CustomHandle != nil {
		if err := h.CustomHandle(ctx, r); err != nil {
			return err
		}
	}

	return h.Handler.Handle(ctx, r)
}

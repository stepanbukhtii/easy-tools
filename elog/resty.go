package elog

import (
	"log/slog"

	"github.com/stepanbukhtii/easy-tools/econtext"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"resty.dev/v3"
)

type RestyLogger struct {
	ServiceName string
}

func (r *RestyLogger) Info(resp *resty.Response, msg string) {
	r.log(resp, nil, msg)
}

func (r *RestyLogger) Error(resp *resty.Response, err error, msg string) {
	r.log(resp, err, msg)
}

func (r *RestyLogger) log(resp *resty.Response, err error, msg string) {
	if resp == nil || resp.Request == nil {
		logger := slog.With(slog.String(string(semconv.ServicePeerNameKey), r.ServiceName))
		if err != nil {
			logger.With(Err(err)).Error(msg)
			logger = logger.With(Err(err))
		}
		logger.Info(msg)
		return
	}

	ctx := resp.Request.Context()

	logger := econtext.Logger(ctx).With(
		slog.String(string(semconv.ServicePeerNameKey), r.ServiceName),
		slog.String(string(semconv.URLFullKey), resp.Request.URL),
		slog.String(string(semconv.HTTPRequestMethodKey), resp.Request.Method),
		slog.Int(string(semconv.HTTPResponseStatusCodeKey), resp.StatusCode()),
	)

	if resp.String() != "" {
		logger = logger.With(slog.String(HTTPResponseBodyContent, resp.String()))
	}

	if resp.Request.Body != nil {
		logger = logger.With(slog.Any(HTTPRequestBodyContent, resp.Request.Body))
	}

	if resp.Request.IsTrace {
		logger = logger.With(
			slog.Group("trace_info",
				slog.Bool("is_conn_reused", resp.Request.TraceInfo().IsConnReused),
				slog.Bool("is_conn_was_idle", resp.Request.TraceInfo().IsConnWasIdle),
				slog.Int64("conn_time", resp.Request.TraceInfo().ConnTime.Milliseconds()),
				slog.Int64("total_time", resp.Request.TraceInfo().TotalTime.Milliseconds()),
				slog.Int64("dns_lookup", resp.Request.TraceInfo().DNSLookup.Milliseconds()),
			),
		)
	}

	if err != nil {
		logger.With(Err(err)).ErrorContext(ctx, msg)
		return
	}

	logger.InfoContext(ctx, msg)
}

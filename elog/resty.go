package elog

import (
	"log/slog"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"resty.dev/v3"
)

//type Resty struct {
//	Service        string
//	Req            []byte
//	Err            error
//	Resp           *resty.Response
//	AdditionalInfo interface{}
//	HideTraceInfo  bool
//	SkipBody       bool
//}
//
//func (l *Resty) Send(ctx context.Context, msg string) {
//	if !restyDebug {
//		return
//	}
//
//	var logEvent *zerolog.Event
//	if l.Err != nil {
//		logEvent = log.Ctx(ctx).Error().Err(l.Err)
//	} else {
//		logEvent = log.Ctx(ctx).Info()
//	}
//
//	if l.Req != nil {
//		logEvent.Bytes(FieldRequestBodyContent, l.Req)
//	}
//
//	if l.Resp == nil {
//		logEvent.Send()
//		return
//	}
//
//	if !l.SkipBody {
//		logEvent.Bytes(FieldResponseBodyContent, l.Resp.Body())
//	}
//
//	logEvent.Int(FieldResponseStatusCode, l.Resp.StatusCode())
//
//	if !l.HideTraceInfo {
//		logEvent.
//			Dict(
//				"trace_info", zerolog.Dict().
//					Bool("is_conn_reused", l.Resp.Request.TraceInfo().IsConnReused).
//					Bool("is_conn_was_idle", l.Resp.Request.TraceInfo().IsConnWasIdle).
//					Int64("conn_time", l.Resp.Request.TraceInfo().ConnTime.Milliseconds()).
//					Int64("total_time", l.Resp.Request.TraceInfo().TotalTime.Milliseconds()).
//					Int64("dns_lookup", l.Resp.Request.TraceInfo().DNSLookup.Milliseconds()),
//			)
//	}
//
//	logEvent.
//		Str(FieldServiceTargetName, l.Service).
//		Str(FieldRequestMethod, l.Resp.Request.Method).
//		Str(FieldURLFull, l.Resp.Request.URL)
//
//	if len(msg) > 0 {
//		logEvent.Msg(msg)
//	} else {
//		logEvent.Send()
//	}
//}

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
		slog.With(Err(err), slog.String(ServiceTargetName, r.ServiceName)).Error(msg)
		return
	}

	ctx := resp.Request.Context()

	logger := econtext.Logger(ctx)
	if !logger.Enabled(ctx, slog.LevelInfo) {
		logger = slog.Default()
	}

	logger = slog.With(
		slog.String(ServiceTargetName, r.ServiceName),
		slog.String(URLOrigin, resp.Request.URL),
		slog.String(HTTPRequestMethod, resp.Request.Method),
		slog.Int(HTTPResponseStatusCode, resp.StatusCode()),
	)

	if resp.String() != "" {
		logger = slog.With(slog.String(HTTPResponseBodyContent, resp.String()))
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

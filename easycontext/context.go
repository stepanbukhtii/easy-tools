package easycontext

import "context"

type traceIDKey struct{}

func AddTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func GetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return "", false
	}
	return traceID, true
}

type logKey struct{}

func InitLogFields(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, map[string]string{})
}

func AddLogField(ctx context.Context, key, value string) context.Context {
	if fields, ok := ctx.Value(logKey{}).(map[string]string); ok {
		fields[key] = value
		return ctx
	}
	return context.WithValue(ctx, logKey{}, map[string]string{key: value})
}

func CopyLogFields(ctx context.Context, key, value string) context.Context {
	if fields, ok := ctx.Value(logKey{}).(map[string]string); ok {
		fields[key] = value
		return context.WithValue(ctx, logKey{}, fields)
	}
	return context.WithValue(ctx, logKey{}, map[string]string{key: value})
}

func GetLogFields(ctx context.Context) (map[string]string, bool) {
	fields, ok := ctx.Value(logKey{}).(map[string]string)
	if !ok {
		return nil, false
	}
	return fields, true
}

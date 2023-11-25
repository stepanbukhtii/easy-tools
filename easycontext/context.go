package easycontext

import "context"

type (
	contextKey[T any] struct{}
	traceIDKey        struct{}
	localeKey         struct{}
	ipAddrKey         struct{}
	subjectKey        struct{}
	rolesKey          struct{}
	skipLoggerKey     struct{}
	logKey            struct{}
)

func WithValue[T any](ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, contextKey[T]{}, value)
}

func Value[T1, T2 any](ctx context.Context) T2 {
	var value T2
	value, _ = ctx.Value(contextKey[T1]{}).(T2)
	return value
}

func GetDefault[T1, T2 any](ctx context.Context, defaultValue T2) T2 {
	var value T2
	var ok bool

	value, ok = ctx.Value(contextKey[T1]{}).(T2)
	if !ok {
		return defaultValue
	}

	return value
}

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func TraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return ""
	}
	return traceID
}

func SetLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeKey{}, locale)
}

func Locale(ctx context.Context) string {
	locale, ok := ctx.Value(localeKey{}).(string)
	if !ok {
		return ""
	}
	return locale
}

func LocaleDefault(ctx context.Context, defaultValue string) string {
	if locale := Locale(ctx); locale != "" {
		return locale
	}
	return defaultValue
}

func SetIPAddress(ctx context.Context, ipAddress string) context.Context {
	return context.WithValue(ctx, ipAddrKey{}, ipAddress)
}

func IPAddress(ctx context.Context) string {
	ippAddress, ok := ctx.Value(ipAddrKey{}).(string)
	if !ok {
		return ""
	}
	return ippAddress
}

func SetSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subjectKey{}, subject)
}

func Subject(ctx context.Context) string {
	subject, ok := ctx.Value(subjectKey{}).(string)
	if !ok {
		return ""
	}
	return subject
}

func SetRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesKey{}, roles)
}

func Roles(ctx context.Context) []string {
	roles, ok := ctx.Value(rolesKey{}).([]string)
	if !ok {
		return nil
	}
	return roles
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

func LogFields(ctx context.Context) map[string]string {
	fields, ok := ctx.Value(logKey{}).(map[string]string)
	if !ok {
		return nil
	}
	return fields
}

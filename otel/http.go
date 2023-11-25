package otel

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewHTTPClient(serviceName string) *http.Client {
	return NewHTTPClientTransport(serviceName, http.DefaultTransport)
}

func NewHTTPClientTransport(serviceName string, transport http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(
			transport,
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("%s: %s %s", serviceName, r.Method, r.URL.String())
			}),
		),
	}
}

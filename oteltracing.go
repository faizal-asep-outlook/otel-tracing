package otelTracing

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type otelTracing struct {
}

func (t *otelTracing) MiddlewareGinTrace() gin.HandlerFunc {
	return MiddlewareGinTrace()
}

func (t *otelTracing) MiddlewareLogger() gin.HandlerFunc {
	return MiddlewareLogger()
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func (t *otelTracing) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return TraceStart(ctx, name)
}

func (t *otelTracing) ShutDown(ctx context.Context) error {
	return ShutDown(ctx)
}

func (t *otelTracing) LogTrace(ctx context.Context, args ...interface{}) {
	LogTrace(ctx, args...)
}

func (t *otelTracing) LogDebug(ctx context.Context, args ...interface{}) {
	LogDebug(ctx, args...)
}

func (t *otelTracing) LogPrint(ctx context.Context, args ...interface{}) {
	LogPrint(ctx, args...)
}

func (t *otelTracing) LogInfo(ctx context.Context, args ...interface{}) {
	LogInfo(ctx, args...)
}

func (t *otelTracing) LogWarn(ctx context.Context, args ...interface{}) {
	LogWarn(ctx, args...)
}

func (t *otelTracing) LogError(ctx context.Context, args ...interface{}) {
	LogError(ctx, args...)
}

func (t *otelTracing) LogFatal(ctx context.Context, args ...interface{}) {
	LogFatal(ctx, args...)
}

func (t *otelTracing) LogPanic(ctx context.Context, args ...interface{}) {
	LogPanic(ctx, args...)
}

func (t *otelTracing) HttpDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	return HttpDo(ctx, req)
}

func (t *otelTracing) HttpGet(ctx context.Context, url string) (*http.Response, error) {
	return HttpGet(ctx, url)
}

func (t *otelTracing) HttpPost(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	return HttpPost(ctx, url, contentType, body)
}

func (t *otelTracing) HttpPostForm(ctx context.Context, url string, data url.Values) (*http.Response, error) {
	return HttpPostForm(ctx, url, data)
}

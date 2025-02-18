package otelTracing

import (
	"context"

	"github.com/gin-gonic/gin"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type otelTracing struct {
}

func (t *otelTracing) MiddlewareGinTrace() gin.HandlerFunc {
	return MiddlewareGinTrace()
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

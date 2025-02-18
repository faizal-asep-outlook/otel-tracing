package otelTracing

import (
	"context"

	"github.com/gin-gonic/gin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type otelTracing struct {
	tp     *sdktrace.TracerProvider
	tracer oteltrace.Tracer
}

func (t *otelTracing) MiddlewareGinTrace() gin.HandlerFunc {
	return MiddlewareGinTrace()
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func (t *otelTracing) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return TraceStart(ctx, name)
}

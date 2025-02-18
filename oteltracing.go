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

package otelTracing

import (
	"context"
	"fmt"

	"github.com/faizal-asep-outlook/otel-tracing/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var Tracer oteltrace.Tracer
var tracerprovider *sdktrace.TracerProvider

type OtelTracing interface {
	MiddlewareGinTrace() gin.HandlerFunc
	TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span)
	ShutDown(ctx context.Context) error
}

func InitTracer() (OtelTracing, error) {
	ctx := context.Background()
	config, err := config.NewConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	rp, err := newResource(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp, err := newTracerProvider(ctx, config, rp)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.Tracer("gin-server")
	tracerprovider = tp
	Tracer = tp.Tracer(config.ServiceName)
	return &otelTracing{}, nil
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	//nolint: spancheck
	return Tracer.Start(ctx, name)
}

// MiddlewareTrace ginMiddleware.
func MiddlewareGinTrace() gin.HandlerFunc {
	Propagators := otel.GetTextMapPropagator()
	return func(c *gin.Context) {

		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		ctx := Propagators.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.HTTPRoute(c.FullPath())),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := Tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		span.SetStatus(_serverStatus(status))
		if status > 0 {
			span.SetAttributes(semconv.HTTPStatusCode(status))
		}
		if len(c.Errors) > 0 {
			span.SetStatus(codes.Error, c.Errors.String())
			for _, err := range c.Errors {
				span.RecordError(err.Err)
			}
		}
	}
}

func ShutDown(ctx context.Context) error {
	return tracerprovider.Shutdown(ctx)
}

func _serverStatus(code int) (codes.Code, string) {
	if code < 100 || code >= 600 {
		return codes.Error, fmt.Sprintf("Invalid HTTP status code %d", code)
	}
	if code >= 500 {
		return codes.Error, ""
	}
	return codes.Unset, ""
}

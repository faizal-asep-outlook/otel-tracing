package otelTracing

import (
	"context"
	"fmt"

	"github.com/faizal-asep-outlook/otel-tracing/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var Tracer oteltrace.Tracer
var tracerprovider *sdktrace.TracerProvider
var logprovider *sdklog.LoggerProvider

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

	lp, err := newLoggerProvider(ctx, config, rp)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create an *otellogrus.Hook and use it in your application.
	hook := otellogrus.NewHook(config.ServiceName, otellogrus.WithLoggerProvider(lp))
	// Set the newly created hook as a global logrus hook
	logrus.AddHook(hook)

	tp, err := newTracerProvider(ctx, config, rp)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.Tracer("gin-server")

	logprovider = lp
	tracerprovider = tp
	Tracer = tp.Tracer(config.ServiceName)
	return &otelTracing{}, nil
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	//nolint: spancheck
	return Tracer.Start(ctx, name)
}

func ShutDown(ctx context.Context) (err error) {
	if err = logprovider.Shutdown(ctx); err != nil {
		return
	}
	if err = tracerprovider.Shutdown(ctx); err != nil {
		return
	}
	return
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

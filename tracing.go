package otelTracing

import (
	"context"
	"fmt"
	"runtime"

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
var loger *logrus.Logger

type noopWriter struct{}

func (noopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type OtelTracing interface {
	MiddlewareGinTrace() gin.HandlerFunc
	TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span)
	ShutDown(ctx context.Context) error
	LogTrace(ctx context.Context, args ...interface{})
	LogDebug(ctx context.Context, args ...interface{})
	LogPrint(ctx context.Context, args ...interface{})
	LogInfo(ctx context.Context, args ...interface{})
	LogWarn(ctx context.Context, args ...interface{})
	LogError(ctx context.Context, args ...interface{})
	LogFatal(ctx context.Context, args ...interface{})
	LogPanic(ctx context.Context, args ...interface{})
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
	log := logrus.New()
	log.AddHook(hook)
	log.SetOutput(&noopWriter{})

	tp, err := newTracerProvider(ctx, config, rp)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.Tracer("gin-server")

	logprovider = lp
	tracerprovider = tp
	loger = log
	Tracer = tp.Tracer(config.ServiceName)
	return &otelTracing{}, nil
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	//nolint: spancheck
	return Tracer.Start(ctx, name)
}

func LogTrace(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Trace(args...)
	} else {
		loger.WithContext(ctx).Trace(args...)
	}
}

func LogDebug(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Debug(args...)
	} else {
		loger.WithContext(ctx).Debug(args...)
	}
}

func LogPrint(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Print(args...)
	} else {
		loger.WithContext(ctx).Print(args...)
	}
}

func LogInfo(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Info(args...)
	} else {
		loger.WithContext(ctx).Info(args...)
	}
}

func LogWarn(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Warn(args...)
	} else {
		loger.WithContext(ctx).Warn(args...)
	}
}

func LogError(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Error(args...)
	} else {
		loger.WithContext(ctx).Error(args...)
	}
}

func LogFatal(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Fatal(args...)
	} else {
		loger.WithContext(ctx).Fatal(args...)
	}
}

func LogPanic(ctx context.Context, args ...interface{}) {
	if _, file, len, ok := runtime.Caller(1); ok {
		loger.WithContext(ctx).WithField("file", fmt.Sprintf("%s(%d)", file, len)).Panic(args...)
	} else {
		loger.WithContext(ctx).Panic(args...)
	}
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

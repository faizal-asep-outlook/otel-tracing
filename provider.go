package otelTracing

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/faizal-asep-outlook/otel-tracing/config"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// newLoggerProvider creates a new logger provider with the OTLP gRPC exporter.
func newLoggerProvider(ctx context.Context, cfg config.Config, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	var (
		exporter sdklog.Exporter
		err      error
	)
	if cfg.OtlpEndpoint == "" {
		exporter, err = stdoutlog.New(stdoutlog.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	} else {
		var secureOption otlploggrpc.Option

		if !cfg.Insecure {
			secureOption = otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		} else {
			secureOption = otlploggrpc.WithInsecure()
		}

		exporter, err = otlploggrpc.New(
			ctx,
			secureOption,
			otlploggrpc.WithEndpoint(cfg.OtlpEndpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
		}
	}

	processor := sdklog.NewBatchProcessor(exporter)
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(res),
	)

	return lp, nil
}

// newTracerProvider creates a new tracer provider with the OTLP gRPC exporter.
func newTracerProvider(ctx context.Context, cfg config.Config, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)
	if cfg.OtlpEndpoint == "" {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	} else {
		var secureOption otlptracegrpc.Option

		if !cfg.Insecure {
			secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		} else {
			secureOption = otlptracegrpc.WithInsecure()
		}
		exporter, err = otlptracegrpc.New(
			ctx,
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.OtlpEndpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
	}

	// Create Resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

// newMeterProvider creates a new meter provider with the OTLP gRPC exporter.
func newMeterProvider(ctx context.Context, cfg config.Config, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	var (
		exporter sdkmetric.Exporter
		err      error
	)

	if cfg.OtlpEndpoint == "" {
		exporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	} else {
		var secureOption otlpmetricgrpc.Option

		if !cfg.Insecure {
			secureOption = otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		} else {
			secureOption = otlpmetricgrpc.WithInsecure()
		}

		exporter, err = otlpmetricgrpc.New(
			ctx,
			secureOption,
			otlpmetricgrpc.WithEndpoint(cfg.OtlpEndpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
		}
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return mp, nil
}

// newResource creates a new OTEL resource with the service name and version.
func newResource(ctx context.Context, cfg config.Config) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String("development"),
			attribute.String("environment", os.Getenv("GO_ENV")),
		),
	)
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

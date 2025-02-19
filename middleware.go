package otelTracing

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

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
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}

		path := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.URL.Path)
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		spanName := path
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

func MiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Start timer
		start := time.Now()

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		path := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.URL.Path)
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		finish := time.Now()
		param := log.Fields{
			// "Request":      c.Request,
			// "Keys":         c.Keys,
			"TimeStamp":    finish.Format("2006-01-02 15:04:05"),
			"Latency":      finish.Sub(start),
			"ClientIP":     c.ClientIP(),
			"Method":       c.Request.Method,
			"StatusCode":   c.Writer.Status(),
			"ErrorMessage": c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"BodySize":     c.Writer.Size(),
		}

		if raw != "" {
			path = path + "?" + raw
		}

		// param.Path = path
		loger.WithContext(c.Request.Context()).WithFields(param).Info(path)
		// fmt.Fprint(out, formatter(param))

	}
}

// MeterRequest is a gin middleware that captures the duration of the request.
func MiddlewareMeter() gin.HandlerFunc {
	// init metric, here we are using histogram for capturing request duration
	histogram, err := MeterInt64Histogram(MetricRequestDurationMillis)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create histogram: %w", err))
	}
	// init metric, here we are using counter for capturing request in flight
	counter, err := MeterInt64UpDownCounter(MetricRequestsInFlight)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create counter: %w", err))
	}

	return func(c *gin.Context) {
		// capture the start time of the request
		startTime := time.Now()

		// define metric attributes

		attrs := metric.WithAttributes(semconv.HTTPRoute(c.FullPath()))

		// increase the number of requests in flight
		counter.Add(c.Request.Context(), 1, attrs)

		// execute next http handler
		c.Next()

		// record the request duration
		duration := time.Since(startTime)
		histogram.Record(
			c.Request.Context(),
			duration.Milliseconds(),
			metric.WithAttributes(
				semconv.HTTPRoute(c.FullPath()),
			),
		)

		// decrease the number of requests in flight
		counter.Add(c.Request.Context(), -1, attrs)
	}
}

// Metric represents a metric that can be collected by the server.
type Metric struct {
	Name        string
	Unit        string
	Description string
}

// MetricRequestDurationMillis is a metric that measures the latency of HTTP requests processed by the server, in milliseconds.
var MetricRequestDurationMillis = Metric{
	Name:        "request_duration_millis",
	Unit:        "ms",
	Description: "Measures the latency of HTTP requests processed by the server, in milliseconds.",
}

// MetricRequestsInFlight is a metric that measures the number of requests currently being processed by the server.
var MetricRequestsInFlight = Metric{
	Name:        "requests_inflight",
	Unit:        "{count}",
	Description: "Measures the number of requests currently being processed by the server.",
}

// MeterInt64Histogram creates a new int64 histogram metric.
func MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error) {
	histogram, err := meter.Int64Histogram(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create histogram: %w", err)
	}

	return histogram, nil
}

// MeterInt64UpDownCounter creates a new int64 up down counter metric.
func MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error) {
	counter, err := meter.Int64UpDownCounter(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create counter: %w", err)
	}

	return counter, nil
}

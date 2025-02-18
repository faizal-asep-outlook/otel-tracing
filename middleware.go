package otelTracing

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

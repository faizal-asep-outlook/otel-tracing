# otel-tracing

## environment variables
```
OTEL_TRACING_OTLP_ENDPOINT=127.0.0.1:4317
OTEL_TRACING_SERVICE_NAME=service
OTEL_TRACING_SERVICE_VERSION=1.0.0
OTEL_TRACING_INSECURE_MODE=true
```

## sample
```
import (
	"context"

	"net/http"

	trace "github.com/faizal-asep-outlook/otel-tracing"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	_, err := trace.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.New()
	r.Use(
		trace.MiddlewareGinTrace(),
		trace.MiddlewareLogger(),
	)

	r.GET("/ping", func(c *gin.Context) {
		_, span := trace.TraceStart(c.Request.Context(), "ping process")
		defer span.End()
		log.WithContext(c.Request.Context()).Info("testing")
		c.String(http.StatusOK, "pong")
	})

	r.Run(":8080")
	trace.ShutDown(context.Background())
}
```
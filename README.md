# otel-tracing
```
import (
	"log"
	"net/http"

	trace "github.com/faizal-asep-outlook/otel-tracing"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := trace.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.New()
	r.Use(trace.MiddlewareGinTrace())

	r.GET("/ping", func(c *gin.Context) {
		_, span := trace.TraceStart(c.Request.Context(), "ping process")
		defer span.End()
		c.String(http.StatusOK, "pong")
	})

	r.Run(":8080")
}
```
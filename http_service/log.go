package httpservice

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLogger logs a gin HTTP request in JSON format. Allows to set the
// logger for testing purposes.
func StructuredLogger() gin.HandlerFunc {
    return func(c *gin.Context) {

        start := time.Now() // Start timer
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
		trace_id := c.Request.Header.Get("X-Trace-Id")
		user_id := c.Request.Header.Get("X-User-Id")
        // Process request
        c.Next()

        // Fill the params
        param := gin.LogFormatterParams{}

        param.TimeStamp = time.Now() // Stop timer
        param.Latency = param.TimeStamp.Sub(start)
        if param.Latency > time.Minute {
            param.Latency = param.Latency.Truncate(time.Second)
        }

        param.ClientIP = c.ClientIP()
        param.Method = c.Request.Method
        param.StatusCode = c.Writer.Status()
        param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
        param.BodySize = c.Writer.Size()
        if raw != "" {
            path = path + "?" + raw
        }
        param.Path = path


        // Log using the params
        var logEvent func(msg string,args ...any)

        if c.Writer.Status() >= 400 {
            logEvent = slog.Error
        } else {
            logEvent = slog.Info
        }
        
        // 可以加入其他的东西，比如trace_id   user 等
        logEvent("http_request",
		slog.String("trace_id", trace_id),
		slog.String("user_id", user_id),
		slog.String("client_id", param.ClientIP),
		slog.String("method", param.Method),
		slog.String("path", param.Path),
		slog.Int("status_code", param.StatusCode),
		slog.Int("body_size", param.BodySize),
		slog.String("latency", param.Latency.String()),
		slog.String("error_message", param.ErrorMessage),
	)
    }
}


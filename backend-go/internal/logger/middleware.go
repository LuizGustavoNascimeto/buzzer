package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GinCloudWatchMiddleware(cwLogger *CloudWatchLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		statusCode := c.Writer.Status()

		level := "INFO"
		if statusCode >= 500 {
			level = "ERROR"
		} else if statusCode >= 400 {
			level = "WARN"
		}

		entry := LogEntry{
			Timestamp:  time.Now(),
			Level:      level,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: statusCode,
			Latency:    time.Since(start),
			ClientIP:   c.ClientIP(),
			Error:      c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		go func() {
			if err := cwLogger.SendLog(entry); err != nil {
				fmt.Printf("falha ao enviar log: %v\n", err)
			}
		}()
	}
}

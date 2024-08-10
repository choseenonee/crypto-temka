package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func (m Middleware) Latency() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()

		//todo: more logs
		m.logger.Info(fmt.Sprintf("handled %v, latency: %v, response status: %v", c.Request.URL.Path, latency, status))
	}
}

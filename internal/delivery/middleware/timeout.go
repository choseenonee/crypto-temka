package middleware

import (
	"context"
	"crypto-temka/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

func (m Middleware) Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*time.Duration(viper.GetInt(config.Timeout)))
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

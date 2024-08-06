package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

const (
	CUserID = "userID"
)

func (m Middleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		auth := c.GetHeader("Authorization")
		if !strings.Contains(auth, "Bearer") {
			m.logger.Info(fmt.Sprintf("unathorized (NO JWT) access at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "no bearer provided in authorization"})
			return
		}

		// WITH Bearer <token>
		jwtToken := strings.Split(auth, " ")[1]

		id, err := m.jwtUtil.Authorize(jwtToken)
		if errors.Is(err, jwt.ErrTokenExpired) {
			m.logger.Info(fmt.Sprintf("token expired at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUpgradeRequired, gin.H{"detail": "JWT expired"})
			return
		}
		if err != nil {
			m.logger.Error(fmt.Sprintf("token parse error at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "error on parsing JWT"})
			return
		}

		c.Set(CUserID, id)

		// BEFORE AND...

		c.Next()

		// ...AFTER THE REQUEST

		latency := time.Since(t)
		status := c.Writer.Status()

		m.logger.Info(fmt.Sprintf("handled %v, latency: %v, response status: %v", c.Request.URL.Path, latency, status))
	}
}

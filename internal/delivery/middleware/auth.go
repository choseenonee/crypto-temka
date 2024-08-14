package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

const (
	CUserID = "userID"
)

func (m Middleware) Authorization(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.Contains(auth, "Bearer") {
			m.logger.Info(fmt.Sprintf("unathorized (NO JWT) access at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "no bearer provided in authorization"})
			return
		}

		// WITH Bearer <token>
		jwtToken := strings.Split(auth, " ")[1]

		id, isAdmin, err := m.jwtUtil.Authorize(jwtToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				m.logger.Info(fmt.Sprintf("token expired at: %v", c.Request.URL.Path))
				c.AbortWithStatusJSON(http.StatusUpgradeRequired, gin.H{"detail": "JWT expired"})
				return
			}
			m.logger.Error(fmt.Sprintf("token parse error at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "error on parsing JWT"})
			return
		}

		if isAdmin != admin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "only admin token for this endpoint"})
			return
		}

		c.Set(CUserID, id)

		// BEFORE AND...

		c.Next()

		// ...AFTER THE REQUEST
	}
}

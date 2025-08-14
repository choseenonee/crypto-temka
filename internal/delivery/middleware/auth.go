package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	CUserID = "userID"
)

func (m Middleware) Authorization(admin bool, getStatus func(id int) (string, error), allowedStatuses []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.Contains(auth, "Bearer") {
			m.logger.Info(fmt.Sprintf("unathorized (NO JWT) access at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no bearer provided in authorization"})
			return
		}

		// WITH Bearer <token>
		jwtToken := strings.Split(auth, " ")[1]

		id, isAdmin, err := m.jwtUtil.Authorize(jwtToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				m.logger.Info(fmt.Sprintf("token expired at: %v", c.Request.URL.Path))
				c.AbortWithStatusJSON(http.StatusUpgradeRequired, gin.H{"message": "JWT expired"})
				return
			}

			m.logger.Error(fmt.Sprintf("token parse error: %v at: %v", err, c.Request.URL.Path))

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("error on parsing JWT: %v", err)})

			return
		}

		if admin {
			if !isAdmin {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "only admin token for this endpoint"})

				return
			}
		} else {
			if isAdmin {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "only user token for this endpoint"})

				return
			}
		}

		if !admin {
			status, err := getStatus(id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("database invalid: %v", err.Error())})
				return
			}

			if !slices.Contains(allowedStatuses, status) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("your account status: %v, "+
					"allowed statuses: %v", status, allowedStatuses)})
				return
			}
		}

		c.Set(CUserID, id)

		// BEFORE AND...

		c.Next()

		// ...AFTER THE REQUEST
	}
}

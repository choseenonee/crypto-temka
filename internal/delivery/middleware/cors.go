package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (m Middleware) CORS() gin.HandlerFunc {
	//return func(c *gin.Context) {
	//	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	//	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	//
	//	if c.Request.Method == "OPTIONS" {
	//		c.AbortWithStatus(http.StatusNoContent)
	//		return
	//	}
	//
	//	c.Next()
	//}

	return cors.New(cors.Config{
		AllowMethods: []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
			"accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		//AllowOrigins:     []string{"https://admin.ai-trade.net", "http://localhost:3000"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 1 * time.Hour,
	})
}

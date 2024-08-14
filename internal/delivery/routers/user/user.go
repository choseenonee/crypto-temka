package user

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine, userService service.User, mdw middleware.Middleware) *gin.RouterGroup {
	router := r.Group("/user")

	router.Use(mdw.Authorization(false, userService.GetStatus))

	handler := handlers.InitUserHandler(userService)

	router.GET("", handler.GetMe)

	return router
}

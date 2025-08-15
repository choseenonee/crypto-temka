package user

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine, userService service.User, mdw middleware.Middleware) *gin.RouterGroup {
	verifiedUserRouter := r.Group("/user")

	verifiedUserRouter.Use(mdw.Authorization(false, userService.GetStatus, []string{"verified", "declined", "opened", "pending"}))

	handler := handlers.InitUserHandler(userService)

	verifiedUserRouter.GET("", handler.GetMe)

	notVerifiedUserRouter := r.Group("/user/properties")

	notVerifiedUserRouter.Use(mdw.Authorization(false, userService.GetStatus, []string{"declined", "opened", "pending"}))

	notVerifiedUserRouter.PUT("", handler.UpdateProperties)

	return nil
}

package user

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterUserRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, mdw middleware.Middleware) *gin.RouterGroup {
	router := r.Group("/user")

	router.Use(mdw.Authorization(false))

	userRepo := repository.InitUser(db)

	jwt := auth.InitJWTUtil()

	userService := service.InitUser(userRepo, jwt, logger)

	handler := handlers.InitUserHandler(userService)

	router.GET("", handler.GetMe)

	return router
}

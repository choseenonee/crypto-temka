package admin

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminUserRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/user")

	userRepo := repository.InitUser(db)

	jwt := auth.InitJWTUtil()

	userService := service.InitUser(userRepo, jwt, logger)

	handler := handlers.InitUserHandler(userService)

	router.GET("/all", handler.GetAll)
	router.PUT("/status", handler.UpdateStatus)

	return router
}

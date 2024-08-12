package routers

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware/auth"
	userrepo "crypto-temka/internal/repository/user"
	userserv "crypto-temka/internal/service/user"
	"crypto-temka/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterUserRouter(r *gin.Engine, db *sqlx.DB, jwtUtil auth.JWTUtil, logger *log.Logs) *gin.RouterGroup {
	userRouter := r.Group("/user")

	userRepo := userrepo.InitUserRepo(db)

	userService := userserv.InitUserService(userRepo, logger, jwtUtil)

	userHandler := handlers.InitUserHandler(userService)

	userRouter.POST("/create", userHandler.Create)
	userRouter.POST("/login", userHandler.Login)
	userRouter.GET("/:id", userHandler.Get)
	userRouter.DELETE("/:id", userHandler.Delete)

	userRouter.PUT("/photo/add", userHandler.AddPhoto)
	userRouter.GET("/photo/:id", userHandler.GetPhoto)

	userRouter.PUT("/status", userHandler.SetStatus)

	return userRouter
}

package public

import (
	"crypto-temka/internal/delivery/handlers"
	"crypto-temka/internal/delivery/middleware/auth"
	"crypto-temka/internal/repository"
	"crypto-temka/internal/service"
	"crypto-temka/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterPublicUserRouter(r *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) *gin.RouterGroup {
	router := r.Group("/user")

	userRepo := repository.InitUser(db)

	jwt := auth.InitJWTUtil()

	walletRepo := repository.InitWallet(db)

	withdrawRepo := repository.InitWithdraw(db)

	userRateRepo := repository.InitUsersRate(db)

	userService := service.InitUser(userRepo, walletRepo, withdrawRepo, userRateRepo, jwt, logger)

	handler := handlers.InitUserHandler(userService)

	router.POST("", handler.Register)
	router.PUT("", handler.Auth)

	return router
}
